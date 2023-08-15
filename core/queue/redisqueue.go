package queue

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"

	redis "github.com/redis/go-redis/v9"
	"github.com/uaxe/infra/core/threading"
)

type IQueue interface {
	Start() error
	Push(rawData []byte) error
	Length() (int64, error)
	Close() error
}

var _ IQueue = (*RedisQueue)(nil)

type RedisQueue struct {
	IQueue
	ctx                  context.Context
	cancel               context.CancelFunc
	qname, stream, group string
	consumerCount        int
	wakeupChan           chan string
	locked               *sync.Map
	client               *redis.Client
	work                 QWorker
}

const (
	QUEUE_STREAM_PREFIX   = "%s:stream"
	QUEUE_GROUP_PREFIX    = "%s:group"
	QUEUE_CONSUMER_PREFIX = "%s:consumer:%d"
	STREAM_START          = "0"
	STREAM_CURRENT        = ">"
)

type RedisQueueOption func(*RedisQueue) error

func (f RedisQueueOption) apply(opts *RedisQueue) error {
	return f(opts)
}

func SetNumConsumer(count int) RedisQueueOption {
	return func(q *RedisQueue) (err error) {
		q.consumerCount = count
		return
	}
}

type QWorker func(channelid string, rawData []byte) error

func NewRedisQueue(ctx context.Context, qname string, client *redis.Client, work QWorker, opts ...RedisQueueOption) (*RedisQueue, error) {
	newCtx, cancel := context.WithCancel(ctx)
	q := &RedisQueue{
		ctx:           newCtx,
		cancel:        cancel,
		qname:         qname,
		stream:        fmt.Sprintf(QUEUE_STREAM_PREFIX, qname),
		group:         fmt.Sprintf(QUEUE_GROUP_PREFIX, qname),
		client:        client,
		consumerCount: runtime.NumCPU() << 1,
		wakeupChan:    make(chan string, 1000),
		locked:        &sync.Map{},
		work:          work,
	}
	for _, opt := range opts {
		if err := opt.apply(q); err != nil {
			return nil, err
		}
	}
	return q, nil
}

func (q *RedisQueue) Start() error {
	if q.work == nil {
		return nil
	}
	threading.GoSafe(func() {
		// 消费历史消息
		length, err := q.Length()
		if err != nil {
			log.Println(err)
			return
		}
		if length > 0 {
			for i := 0; i < q.consumerCount; i++ {
				_, loaded := q.locked.LoadOrStore(i, "1")
				if !loaded {
					go func(idx int) {
						defer func() {
							q.locked.Delete(idx)
						}()
						consumer := fmt.Sprintf(QUEUE_CONSUMER_PREFIX, q.qname, idx)
						q.consumeHistory(consumer, STREAM_START)
					}(i)
				}
			}
		}
	})
	threading.GoSafe(func() {
		q.startCore()
	})
	// 唤醒处理
	q.wakeupChan <- STREAM_CURRENT
	// 创建一个消费组,消费新产生的消息
	q.client.XGroupCreateMkStream(q.ctx, q.stream, q.group, "0").Result()
	return nil
}

func (q *RedisQueue) consumeHistory(consumer, start string) {
	limit := 10
	args := &redis.XReadArgs{
		Streams: []string{q.stream, start},
		Count:   int64(limit),
		Block:   100 * time.Millisecond,
	}
	for {
		results, err := q.client.XRead(q.ctx, args).Result()
		if err != nil && err != redis.Nil {
			log.Println(err)
			break
		}
		for _, re := range results {
			for _, msg := range re.Messages {
				if len(msg.Values) == 0 {
					continue
				}
				q.consumeOne(consumer, msg)
			}
		}
		if len(results) < limit {
			break
		}
	}
}

func (q *RedisQueue) Close() error {
	q.cancel()
	time.Sleep(5 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	for i := 0; i < q.consumerCount; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			consumer := fmt.Sprintf(QUEUE_CONSUMER_PREFIX, q.qname, idx)
			pendings, err := q.client.XPendingExt(ctx, &redis.XPendingExtArgs{
				Stream:   q.stream,
				Group:    q.group,
				Consumer: consumer,
				Count:    100,
				Start:    "-",
				End:      "+",
			}).Result()
			if err != nil {
				log.Printf("redis xpending err: %s", err)
				return
			}
			if len(pendings) == 0 {
				// log.Printf("not pendings: %s,%s,%s", q.stream, q.group, consumer)
				return
			}
			for _, p := range pendings {
				msgID := p.ID
				msglist, err := q.client.XRangeN(ctx, q.stream, msgID, "+", 1).Result()
				if err != nil {
					log.Printf("redis xrange err: %s,%s", err, msgID)
					continue
				}
				for _, msg := range msglist {
					if _, err = q.client.XAdd(ctx, &redis.XAddArgs{
						Stream: q.stream,
						ID:     "*",
						Values: msg.Values,
					}).Result(); err != nil {
						log.Printf("redis xadd err: %+v", err)
						continue
					}
				}
			}
			if _, err := q.client.XGroupDelConsumer(ctx, q.stream, q.group, consumer).Result(); err != nil {
				log.Printf("redis xgroup delconsumer err: %+v", err)
			}
		}(i)
	}
	wg.Wait()
	return nil
}

func (q *RedisQueue) consume(consumer, start string) {
	limit := 50
	args := &redis.XReadGroupArgs{
		Streams:  []string{q.stream, start},
		Group:    q.group,
		Consumer: consumer,
		Count:    int64(limit),
		Block:    100 * time.Millisecond,
		NoAck:    true,
	}
	for {
		results, err := q.client.XReadGroup(q.ctx, args).Result()
		if err != nil && err != redis.Nil {
			log.Println(err)
			break
		}
		for _, re := range results {
			for _, msg := range re.Messages {
				if len(msg.Values) == 0 {
					continue
				}
				q.consumeOne(args.Consumer, msg)
			}
		}
		if len(results) < limit {
			break
		}
	}
}

func (q *RedisQueue) startCore() {
	for {
		select {
		case <-q.ctx.Done():
			q.Close()
			return
		case start := <-q.wakeupChan:
			// 叫醒所有消费者进行处理。
			for i := 0; i < q.consumerCount; i++ {
				_, loaded := q.locked.LoadOrStore(i, "1")
				if !loaded {
					go func(idx int) {
						defer func() {
							q.locked.Delete(idx)
						}()
						consumer := fmt.Sprintf(QUEUE_CONSUMER_PREFIX, q.qname, idx)
						q.consume(consumer, start)
					}(i)
				}
			}
		}
	}
}

func (q *RedisQueue) consumeOne(consumer string, msg redis.XMessage) {
	threading.RunSafe(func() {
		startTime := time.Now()
		defer func() {
			cost := time.Now().Sub(startTime)
			if rand.Intn(1000) == 0 && cost.Milliseconds() > 1000 {
				log.Printf("%s,%d,%s", consumer, cost.Milliseconds(), msg.ID)
			}
		}()
		if q.work != nil {
			value, ok := msg.Values["q"]
			if !ok {
				log.Printf("%s,%s,%+v", consumer, msg.ID, msg.Values)
				return
			}
			rawdata, ok := value.(string)
			if !ok {
				log.Printf("%s,%s,%+v", consumer, msg.ID, msg.Values)
				return
			}
			raw, err := base64.StdEncoding.DecodeString(rawdata)
			if err != nil {
				log.Printf("%s,%s,%+v", consumer, msg.ID, msg.Values)
				return
			}
			if err := q.work(msg.ID, raw); err != nil {
				log.Printf("%s,%+v", err, msg)
			} else {
				q.client.XDel(q.ctx, q.stream, msg.ID)
			}
		}
	})
}

// 获取队列长度
func (q *RedisQueue) Length() (int64, error) {
	return q.client.XLen(q.ctx, q.stream).Result()
}

// 添加队列
func (q *RedisQueue) Push(rawData []byte) error {
	strData := base64.StdEncoding.EncodeToString(rawData)
	values := map[string]interface{}{"q": strData}
	_, err := q.client.XAdd(q.ctx, &redis.XAddArgs{Stream: q.stream, ID: "*", Values: values}).Result()
	if err != nil {
		return err
	}
	//全体唤醒
	q.wakeupChan <- STREAM_CURRENT
	return nil
}
