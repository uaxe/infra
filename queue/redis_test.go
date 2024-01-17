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

	master, slave := hs.FindForClient("100001", nil)
	t.Logf("TestRedisHash|FindForClient(100001, nil)|%s|%s\n", master.Hostport, slave.Hostport)

	sc := hs.ShardNum()
	if sc != 2 {
		t.Fail()
	}
	t.Logf("ShardNum|%d\n", sc)
}
