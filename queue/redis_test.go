package queue_test

import (
	"testing"

	"github.com/uaxe/infra/queue"
)

func TestRedisHash(t *testing.T) {

	options := queue.RedisOptions{
		SlaveOpen:   true,
		ClusterName: "session",
		BasicPort:   8080,
		Username:    "root",
		Password:    "",
		ShardNum:    2,
		ShardSeed:   2,
		MaxIdleConn: 2,
		MaxOpenConn: 4}

	hs := queue.NewRedisShard(options)

	sc := hs.ShardNum()
	if sc != 2 {
		t.Fail()
	}
}
