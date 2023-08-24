package queue

import (
	"context"
	"io"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/redis/go-redis/v9"
)

type (
	RedisQueueManager struct {
		opts *options
	}

	options struct {
		ctx          context.Context
		cancel       context.CancelFunc
		output       io.Writer
		queuenames   []string
		name2Queue   *sync.Map
		name2Redis   *sync.Map
		name2Worker  *sync.Map
		name2Options *sync.Map
	}

	optionFunc func(*options) error
)

func (f optionFunc) apply(opts *options) error {
	return f(opts)
}

// 设置打印输出
func SetQueueManagerOuput(output io.Writer) optionFunc {
	return func(opts *options) error {
		opts.output = output
		return nil
	}
}

func newRedisQueueManagerOptions(ctx context.Context) *options {
	newCtx, cancel := context.WithCancel(ctx)
	return &options{
		ctx:          newCtx,
		cancel:       cancel,
		output:       os.Stdout,
		queuenames:   make([]string, 0, 2),
		name2Queue:   &sync.Map{},
		name2Redis:   &sync.Map{},
		name2Worker:  &sync.Map{},
		name2Options: &sync.Map{},
	}
}

func NewRedisQueueManager(ctx context.Context, opts ...optionFunc) (*RedisQueueManager, error) {
	m := &RedisQueueManager{
		opts: newRedisQueueManagerOptions(ctx),
	}
	for _, opt := range opts {
		if err := opt.apply(m.opts); err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (self *RedisQueueManager) Start() error {
	for _, qname := range self.opts.queuenames {
		clinet, ok := self.opts.name2Redis.Load(qname)
		if !ok {
			continue
		}
		redcli, ok := clinet.(*redis.Client)
		if !ok {
			log.Println("redcli Not Found:", qname)
			continue
		}
		var worker QWorker
		qworker, ok := self.opts.name2Worker.Load(qname)
		if ok {
			worker, ok = qworker.(QWorker)
		}
		var opts []RedisQueueOption
		options, ok := self.opts.name2Options.Load(qname)
		if ok {
			opts, ok = options.([]RedisQueueOption)
		}
		redisQueue, err := NewRedisQueue(self.opts.ctx, qname, redcli, worker, opts...)
		if err != nil {
			return err
		}
		self.opts.name2Queue.Store(qname, redisQueue)
		log.Println("Queue Start:", qname)
		err = redisQueue.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *RedisQueueManager) WatchQueue(queueName string, client *redis.Client, onMessage QWorker, opts ...RedisQueueOption) error {
	sort.Strings(self.opts.queuenames)
	idx := sort.SearchStrings(self.opts.queuenames, queueName)
	if idx == len(self.opts.queuenames) || self.opts.queuenames[idx] != queueName {
		self.opts.queuenames = append(self.opts.queuenames, queueName)
	}
	if nil != client {
		self.opts.name2Redis.Store(queueName, client)
	}
	if nil != onMessage {
		self.opts.name2Worker.Store(queueName, onMessage)
	}
	if len(opts) > 0 {
		self.opts.name2Options.Store(queueName, opts)
	}
	return nil
}

func (self *RedisQueueManager) SelectQueue(qname string) (*RedisQueue, bool) {
	v, loaded := self.opts.name2Queue.Load(qname)
	if !loaded {
		log.Println("selectQueue:", qname, "loaded:", loaded, "queuenames")
		return nil, false
	}
	return v.(*RedisQueue), true
}

// 遍历队列信息
func (self *RedisQueueManager) RangeQueue(iterator func(qname string, queue *RedisQueue)) {
	self.opts.name2Queue.Range(func(key, value interface{}) bool {
		iterator(key.(string), value.(*RedisQueue))
		return true
	})
}

// 获取所有被管理的 queues 的 name
func (self *RedisQueueManager) GetManagedQueuenames() []string {
	return self.opts.queuenames
}

func (self *RedisQueueManager) Stop() {
	self.RangeQueue(func(qname string, queue *RedisQueue) {
		queue.Close()
	})
}
