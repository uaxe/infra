package queue

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisNode struct {
	Client   *redis.Client
	Hostport string
}

type redisShardRange struct {
	min     int
	max     int
	shardId int
	Master  []*RedisNode
	Slave   []*RedisNode
}

type RedisShard struct {
	shardSeed   int
	shardNum    int
	hashNum     int
	shardranges []redisShardRange
	options     RedisOptions
}

type RedisOptions struct {
	SlaveOpen   bool
	ShardSeed   int
	ShardNum    int
	ClusterName string
	BasicPort   int
	Username    string
	Password    string
	MaxIdleConn int
	MaxOpenConn int
	SSLOn       bool
}

func (opt *RedisOptions) String() string {
	return fmt.Sprintf("%s:%d", opt.ClusterName, opt.BasicPort)
}

func NewRedisShard(options RedisOptions) *RedisShard {

	hash := options.ShardSeed / options.ShardNum
	options.ClusterName, _ = url.QueryUnescape(options.ClusterName)

	uniqClients := make(map[string] /*addr*/ *redis.Client, 5)

	shardranges := make([]redisShardRange, 0, hash)
	for i := 0; i < options.ShardNum; i++ {

		clusterMDomain := net.JoinHostPort(options.ClusterName, strconv.Itoa(options.BasicPort))
		clusterSDomain := clusterMDomain
		if options.SlaveOpen {
			clusterSDomain = net.JoinHostPort(options.ClusterName, strconv.Itoa(options.BasicPort))
		}

		rcm, ok := uniqClients[clusterMDomain]
		if !ok {
			opt := &redis.Options{
				Addr:            clusterMDomain,
				Password:        options.Password, // no password set
				DB:              0,                // use default DB
				MaxRetries:      3,
				DialTimeout:     20 * time.Second,
				PoolSize:        options.MaxOpenConn,
				MinIdleConns:    options.MaxIdleConn,
				ConnMaxLifetime: 1 * time.Minute,
			}

			if options.SSLOn {
				host := strings.Split(clusterMDomain, ":")[0]
				opt.TLSConfig = &tls.Config{ServerName: host}
			}
			rcm = redis.NewClient(opt)
			uniqClients[clusterMDomain] = rcm
		}

		rcs, ok := uniqClients[clusterSDomain]
		if !ok {
			opt := &redis.Options{
				Addr:            clusterSDomain,
				Password:        options.Password, // no password set
				DB:              0,                // use default DB
				MaxRetries:      3,
				DialTimeout:     20 * time.Second,
				PoolSize:        options.MaxOpenConn,
				MinIdleConns:    options.MaxIdleConn,
				ConnMaxLifetime: 1 * time.Minute,
			}
			if options.SSLOn {
				host := strings.Split(clusterSDomain, ":")[0]
				opt.TLSConfig = &tls.Config{ServerName: host}
			}
			rcs = redis.NewClient(opt)
			uniqClients[clusterSDomain] = rcs
		}

		master := make([]*RedisNode, 0, hash)
		slave := make([]*RedisNode, 0, hash)
		for j := 0; j < hash; j++ {
			master = append(master, &RedisNode{Client: rcm, Hostport: clusterMDomain})
			slave = append(slave, &RedisNode{Client: rcs, Hostport: clusterSDomain})
		}

		shardranges = append(shardranges,
			redisShardRange{
				min:     i * hash,
				max:     (i + 1) * hash,
				shardId: i,
				Master:  master,
				Slave:   slave})
	}

	return &RedisShard{
		shardSeed:   options.ShardSeed,
		shardNum:    options.ShardNum,
		hashNum:     hash,
		shardranges: shardranges,
		options:     options}
}

func (s *RedisShard) FindForClient(key string,
	hashKeyFunc func(key string) int) (*RedisNode, *RedisNode) {
	if nil == hashKeyFunc {
		hashKeyFunc = defaultHashKeyFunc
	}
	i := hashKeyFunc(key) % s.options.ShardSeed

	for _, v := range s.shardranges {
		if v.min <= i && v.max > i {
			master := v.Master[i%s.hashNum]
			slave := v.Slave[i%s.hashNum]
			return master, slave
		}
	}
	return nil, nil
}

func (s *RedisShard) ShardNum() int {
	return s.shardNum
}

func (s *RedisShard) ShardSeed() int {
	return s.shardSeed
}

func (s *RedisShard) HashNum() int {
	return s.hashNum
}

func (s *RedisShard) Stop() {
	for _, v := range s.shardranges {
		for _, master := range v.Master {
			_ = master.Client.Close()
		}
		for _, slave := range v.Slave {
			_ = slave.Client.Close()
		}
	}
}

var defaultHashKeyFunc = func(key string) int {
	num := fmt.Sprintf("%x", md5.Sum([]byte(key)))
	if len(num) > 1 {
		num = num[len(num)-2:]
	}
	i, err := strconv.ParseInt(num, 16, 64)
	if err != nil {
		return 0
	}
	return int(i)
}

var hashByTail = func(key string) int {
	num := fmt.Sprintf("%x", md5.Sum([]byte(key)))
	num = num[len(num)-2:]
	i, err := strconv.Atoi(num)
	if err != nil {
		return 0
	}
	return i
}

type notifyItem struct {
	ctx         context.Context
	hashid      string
	notifyTopic string
	notifyChan  chan *struct{}
	key         string
	redisNode   *RedisNode
}

type RedisQueue struct {
	IQueue
	redisInstance *RedisShard
	work          QWorker
	meta          QueueMeta
	ctx           context.Context
	cancel        context.CancelFunc
	wakeupChan    chan any
	notifyItems   []*notifyItem
	topic2Items   map[string] /*notifyTopic*/ *notifyItem
	submitTasks   *sync.Map
}

const (
	KeyQueuePrefix       = "_%s:%s:queue_"
	KeyNotifyTopicPrefix = "_%s:%s:topic_"
)

func NewRedisQueue(queueMeta QueueMeta, redisInstance *RedisShard,
	work func(channelid string, raw []byte) error) *RedisQueue {

	ctx, cancel := context.WithCancel(queueMeta.Ctx)
	self := &RedisQueue{
		ctx:           ctx,
		cancel:        cancel,
		meta:          queueMeta,
		redisInstance: redisInstance,
		work:          work,
		submitTasks:   &sync.Map{},
		topic2Items:   make(map[string] /*notifyTopic*/ *notifyItem, queueMeta.HashSize),
		wakeupChan:    make(chan any, 1000),
	}

	for i := 0; i < self.meta.HashSize; i++ {
		m, _ := self.redisInstance.FindForClient(strconv.Itoa(i), func(key string) int {
			v, _ := strconv.Atoi(key)
			return v
		})

		item := &notifyItem{
			notifyChan: make(chan *struct{}, 1),
			redisNode:  m,
			hashid:     strconv.Itoa(i),
			ctx:        ctx,
		}

		qname := fmt.Sprintf(KeyQueuePrefix, self.meta.QueueNamePrefix, item.hashid)
		topic := fmt.Sprintf(KeyNotifyTopicPrefix, self.meta.QueueNamePrefix, item.hashid)

		if queueMeta.TopicMode {
			item.key = topic
			item.notifyTopic = topic
		} else {
			item.notifyTopic = topic
			item.key = qname
		}

		self.topic2Items[item.notifyTopic] = item
		self.notifyItems = append(self.notifyItems, item)
	}

	if queueMeta.TopicMode {
		m, _ := self.redisInstance.FindForClient("0", func(key string) int {
			v, _ := strconv.Atoi(key)
			return v
		})

		item := &notifyItem{
			notifyChan: make(chan *struct{}, 1),
			redisNode:  m,
			hashid:     "0",
			ctx:        ctx,
		}

		item.key = self.meta.QueueNamePrefix
		item.notifyTopic = self.meta.QueueNamePrefix

		self.topic2Items[item.notifyTopic] = item
		self.notifyItems = append(self.notifyItems, item)
	}

	return self
}

func (q *RedisQueue) Start() error {
	redisNodes := make(map[string]*RedisNode)
	subscribes := make(map[string][]string)
	for i := range q.notifyItems {
		item := q.notifyItems[i]
		hostport := item.redisNode.Hostport
		if _, ok := redisNodes[hostport]; !ok {
			redisNodes[hostport] = item.redisNode
		}
		v, ok := subscribes[hostport]
		if !ok {
			v = make([]string, 0, 2)
		}
		subscribes[hostport] = append(v, item.notifyTopic)
	}

	if q.meta.TopicMode {
		subChannels := make([]*redis.PubSub, 0, len(subscribes))
		for node, topics := range subscribes {
			if redisnode, ok := redisNodes[node]; ok {
				sub := redisnode.Client.Subscribe(q.ctx, topics...)
				subChannels = append(subChannels, sub)
			} else {
				panic(fmt.Errorf("no reidsNode [%s]", node))
			}
		}
		q.startTopics(q.work, subChannels...)
	} else {
		wakeupQueuePop := func(topic string, _ []byte) error {
			item, ok := q.topic2Items[topic]
			if ok {
				select {
				case item.notifyChan <- nil:
					select {
					case q.wakeupChan <- nil:
					default:
					}
				default:
				}
			}
			return nil
		}

		notifyitems := make([]*notifyItem, 0, len(subscribes))
		subChannels := make([]*redis.PubSub, 0, len(subscribes))
		for node, topics := range subscribes {
			for _, topic := range topics {
				if _, ok := q.topic2Items[topic]; ok {
					notifyitems = append(notifyitems, q.topic2Items[topic])
				}
			}
			if redisnode, ok := redisNodes[node]; ok {
				sub := redisnode.Client.Subscribe(q.ctx, topics...)
				subChannels = append(subChannels, sub)
			} else {
				panic(fmt.Errorf("no reidsNode %s", node))
			}
		}
		q.startTopics(wakeupQueuePop, subChannels...)
		q.startCore(notifyitems...)
	}
	return nil
}

func (q *RedisQueue) NotifyAll() {
	for _, item := range q.notifyItems {
		select {
		case item.notifyChan <- nil:
		default:
		}
		fmt.Println("RedisQueue|NotifyAll...", item.key, item.notifyTopic)
	}
	select {
	case q.wakeupChan <- nil:
	default:
	}
}

func (q *RedisQueue) startCore(items ...*notifyItem) {
	go func() {
		for {
			select {
			case <-q.ctx.Done():
				_ = q.Close()
				return
			case <-q.wakeupChan:
				for i := range items {
					item := items[i]
					select {
					case <-q.ctx.Done():
						_ = q.Close()
						return
					case <-item.notifyChan:
						_, loaded := q.submitTasks.LoadOrStore(item.key, 1)
						if !loaded {
							go func() {
								defer func() {
									q.submitTasks.Delete(item.key)
								}()
								err := q.handle0(item)
								if err != nil {
									fmt.Println("redisQueue|notify|handleCore|Fail", err, item.key)
								}
							}()
						}
					default:

					}
				}
			}
		}
	}()
}

func (q *RedisQueue) Push(ctx context.Context, hashid string, raw []byte) (bool, error) {

	idx := hashByTail(hashid) % q.meta.HashSize
	item := q.notifyItems[idx]

	strData := base64.StdEncoding.EncodeToString(raw)

	item.redisNode.Client.RPush(ctx, item.key, strData)

	item.redisNode.Client.Publish(ctx, item.notifyTopic, base64.StdEncoding.EncodeToString([]byte{1}))

	select {
	case item.notifyChan <- nil:
	default:
	}

	select {
	case q.wakeupChan <- nil:
	default:
	}

	return true, nil
}

func (q *RedisQueue) Publish(ctx context.Context, hashid string, raw []byte) {
	idx := hashByTail(hashid) % q.meta.HashSize
	item := q.notifyItems[idx]
	strData := base64.StdEncoding.EncodeToString(raw)
	item.redisNode.Client.Publish(ctx, item.notifyTopic, strData)
}

func (q *RedisQueue) startTopics(onTopic func(channelid string, raw []byte) error, pubsubs ...*redis.PubSub) {

	for i := range pubsubs {
		sub := pubsubs[i]
		subChan := sub.Channel()
		go func() {
			for {
				select {
				case <-q.ctx.Done():
					_ = sub.Close()
					return
				case msg := <-subChan:
					topicChannel := msg.Channel
					raw, err := base64.StdEncoding.DecodeString(msg.Payload)
					if err == nil {
						func() {
							defer func() {
								if err := recover(); err != nil {
									fmt.Println("redisQueue|subscribe|listener", topicChannel, string(debug.Stack()))
								}
							}()
							_ = onTopic(topicChannel, raw)
						}()
					}
				}
			}
		}()
	}
}

var (
	ErrNoData = fmt.Errorf("no data error")
)

func (q *RedisQueue) handle0(item *notifyItem) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("RedisQueue|handle0|Panic", item.key, string(debug.Stack()))
		}
	}()
	for {

		select {
		case <-q.ctx.Done():
			return nil
		default:
			raw, err := item.redisNode.Client.LPop(q.ctx, item.key).Result()
			if err != nil && !errors.Is(err, redis.Nil) {
				fmt.Println("RedisQueue|handle0|LPop|Fail", err, item.key)
				return err
			}

			if errors.Is(err, redis.Nil) || len(raw) == 0 {
				return nil
			}

			if len(raw) > 0 {
				newraw, _ := base64.StdEncoding.DecodeString(raw)
				if q.work != nil {
					now := time.Now()
					if err = q.work(item.key, newraw); err != nil {
						fmt.Println("RedisQueue|handle0|BLPop.work|FAIL", err, item.key)
					}
					cost := time.Since(now)
					if rand.Intn(1000) == 0 && cost.Milliseconds() > 1000 {
						fmt.Println("RedisQueue|handle0|BLPop.work|SLOW", item.key, cost.Milliseconds())
					}
				} else {
					fmt.Println("RedisQueue|handle0|BLPop.work|NoWork", item.key)
				}
			}
		}
	}
}

func (q *RedisQueue) Length() int {

	length := 0
	for _, item := range q.notifyItems {
		l, err := item.redisNode.Client.LLen(q.ctx, item.key).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			fmt.Println("RedisQueue|Length|Fail", err, item.key)
			continue
		}
		length += int(l)
	}

	return length
}

func (q *RedisQueue) QueueURL() string {
	return q.redisInstance.options.String()
}

func (q *RedisQueue) Close() error {
	q.cancel()
	fmt.Println("RedisQueue|Close|SUCC", q.meta.QueueNamePrefix, q.QueueURL())
	return nil
}
