package queue

import "context"

type RedisQ struct {
	IFace
	work   QWorker
	ctx    context.Context
	cancel context.CancelFunc
	wakeup chan any
}
