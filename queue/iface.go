package queue

import (
	"context"
)

type (
	IFace interface {
		Start() error
		Push(ctx context.Context, hashid string, raw []byte) (bool, error)
		Publish(ctx context.Context, hashid string, raw []byte)
		Length(hashid string) (int, error)
		Close(hashid string) error
	}

	QWorker func(channelid string, raw []byte) error
)

const (
	QueuePrefix = "_%s:%s:queue_"
	TopicPrefix = "_%s:%s:topic_"
)
