package queue

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	rediscli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	qname := "test"
	q, err := NewRedisQueue(context.Background(), qname, rediscli,
		func(id string, rawData []byte) error {
			rd := rand.Intn(100)
			time.Sleep(time.Duration(rd) * time.Millisecond)
			fmt.Println(id)
			return nil
		}, SetNumConsumer(1))
	assert.Nil(t, err)
	q.Start()

	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(idx int) {
			q.Push([]byte(fmt.Sprintf("%d", idx)))
			wg.Done()
		}(i)
	}
	wg.Wait()
	time.Sleep(10 * time.Second)
	rediscli.Del(context.Background(), qname)
}

func TestGroup(t *testing.T) {
	rediscli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	stream := "notify_test"
	group := "fast"

	var err error

	rediscli.XGroupCreate(context.Background(), stream, group, "0").Result()

	go func() {
		consumers := []string{"ME", "YOU"}
		var wg sync.WaitGroup
		for _, consumer := range consumers {
			wg.Add(1)
			go func(consumerName string) {
				defer wg.Done()
				args := &redis.XReadGroupArgs{
					Streams:  []string{stream, "0"},
					Group:    group,
					Consumer: consumerName,
					Count:    1,
					Block:    100 * time.Millisecond,
					NoAck:    true,
				}
				for {
					xs, err := rediscli.XReadGroup(context.Background(), args).Result()
					if err != nil && err != redis.Nil {
						t.Logf("ME FAIL %s", err)
						break
					}
					for _, s := range xs {
						for _, m := range s.Messages {
							t.Logf("%s %s,%s,%+v", args.Consumer, s.Stream, m.ID, m.Values)
							// rediscli.XAck(context.Background(), s.Stream, group, m.ID)
						}
					}
				}
			}(consumer)
		}
		wg.Wait()
	}()

	// _, err = rediscli.XAdd(context.Background(), &redis.XAddArgs{
	// 	Stream: stream,
	// 	Values: map[string]interface{}{"id": "12"},
	// }).Result()
	// assert.Nil(t, err)

	count, err := rediscli.XLen(context.Background(), stream).Result()
	assert.Nil(t, err)
	t.Logf("%d", count)

	pendings, err := rediscli.XPending(context.Background(), stream, group).Result()
	assert.Nil(t, err)
	t.Logf("%+v", pendings)

	groups, err := rediscli.XInfoGroups(context.Background(), stream).Result()
	assert.Nil(t, err)
	t.Logf("%+v", groups)

	time.Sleep(5 * time.Second)
}

func TestQueueManager(t *testing.T) {
	rediscli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	manager, err := NewRedisQueueManager(context.Background())
	assert.Nil(t, err)
	manager.WatchQueue("test", rediscli, func(channelid string, rawData []byte) error {
		t.Logf("%s", string(rawData))
		return nil
	})
	manager.WatchQueue("test2", rediscli, func(channelid string, rawData []byte) error {
		t.Logf("%s", string(rawData))
		return nil
	})
	manager.Start()
	manager.RangeQueue(func(qname string, queue *RedisQueue) {
		t.Logf("%s", qname)
	})
	defer manager.Stop()

	queue1, succ := manager.SelectQueue("test")
	assert.True(t, succ)
	queue2, succ := manager.SelectQueue("test2")
	assert.True(t, succ)

	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(idx int) {
			queue1.Push([]byte(fmt.Sprintf("q1:%d", idx)))
			queue2.Push([]byte(fmt.Sprintf("q2:%d", idx)))
			wg.Done()
		}(i)
	}
	wg.Wait()
	time.Sleep(10 * time.Second)
	rediscli.Del(context.Background(), "test", "test2")
}
