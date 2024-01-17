package queue

import (
	"context"
)

type (
	IQueue interface {
		Start() error
		Push(ctx context.Context, hashid string, raw []byte) (bool, error)
		Publish(ctx context.Context, hashid string, raw []byte)
		Length(hashid string) (int, error)
		Close(hashid string) error
	}

	QueueMeta struct {
		Ctx             context.Context
		QueueNamePrefix string
		HashSize        int
		TopicMode       bool
	}

	QWorker func(channelid string, raw []byte) error
)
