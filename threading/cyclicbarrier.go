package threading

import (
	"context"
	"errors"
	"sync"
)

type CyclicBarrier interface {
	// 等待所有的参与者到达，如果被ctx.Done()中断，会返回ErrBrokenBarrier
	Await(ctx context.Context) error

	// 重置循环栅栏到初始化状态。如果当前有等待者，那么它们会返回ErrBrokenBarrier
	Reset()

	// 返回当前等待者的数量
	GetNumberWaiting() int

	// 参与者的数量
	GetParties() int

	// 循环栅栏是否处于中断状态
	IsBroken() bool
}

var (
	ErrBrokenBarrier = errors.New("broken barrier")
)

// round
type round struct {
	count    int           // count of goroutines for this roundtrip
	waitCh   chan struct{} // wait channel for this roundtrip
	brokeCh  chan struct{} // channel for isBroken broadcast
	isBroken bool          // is barrier broken
}

// cyclicBarrier impl CyclicBarrier intf
type cyclicBarrier struct {
	parties       int
	barrierAction func() error

	lock  sync.RWMutex
	round *round
}

// New initializes a new instance of the CyclicBarrier, specifying the number of parties.
func New(parties int) CyclicBarrier {
	if parties <= 0 {
		panic("parties must be positive number")
	}
	return &cyclicBarrier{
		parties: parties,
		lock:    sync.RWMutex{},
		round: &round{
			waitCh:  make(chan struct{}),
			brokeCh: make(chan struct{}),
		},
	}
}

// NewWithAction initializes a new instance of the CyclicBarrier,
// specifying the number of parties and the barrier action.
func NewWithAction(parties int, barrierAction func() error) CyclicBarrier {
	if parties <= 0 {
		panic("parties must be positive number")
	}
	return &cyclicBarrier{
		parties: parties,
		lock:    sync.RWMutex{},
		round: &round{
			waitCh:  make(chan struct{}),
			brokeCh: make(chan struct{}),
		},
		barrierAction: barrierAction,
	}
}

func (b *cyclicBarrier) Await(ctx context.Context) error {
	var (
		ctxDoneCh <-chan struct{}
	)
	if ctx != nil {
		ctxDoneCh = ctx.Done()
	}

	// check if context is done
	select {
	case <-ctxDoneCh:
		return ctx.Err()
	default:
	}

	b.lock.Lock()

	// check if broken
	if b.round.isBroken {
		b.lock.Unlock()
		return ErrBrokenBarrier
	}

	// increment count of waiters
	b.round.count++

	// saving in local variables to prevent race
	waitCh := b.round.waitCh
	brokeCh := b.round.brokeCh
	count := b.round.count

	b.lock.Unlock()

	if count > b.parties {
		panic("CyclicBarrier.Await is called more than count of parties")
	}

	if count < b.parties {
		// wait other parties
		select {
		case <-waitCh:
			return nil
		case <-brokeCh:
			return ErrBrokenBarrier
		case <-ctxDoneCh:
			b.breakBarrier(true)
			return ctx.Err()
		}
	} else {
		// we are last, run the barrier action and reset the barrier
		if b.barrierAction != nil {
			err := b.barrierAction()
			if err != nil {
				b.breakBarrier(true)
				return err
			}
		}
		b.reset(true)
		return nil
	}
}

func (b *cyclicBarrier) breakBarrier(needLock bool) {
	if needLock {
		b.lock.Lock()
		defer b.lock.Unlock()
	}

	if !b.round.isBroken {
		b.round.isBroken = true

		// broadcast
		close(b.round.brokeCh)
	}
}

func (b *cyclicBarrier) reset(safe bool) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if safe {
		// broadcast to pass waiting goroutines
		close(b.round.waitCh)

	} else if b.round.count > 0 {
		b.breakBarrier(false)
	}

	// create new round
	b.round = &round{
		waitCh:  make(chan struct{}),
		brokeCh: make(chan struct{}),
	}
}

func (b *cyclicBarrier) Reset() {
	b.reset(false)
}

func (b *cyclicBarrier) GetParties() int {
	return b.parties
}

func (b *cyclicBarrier) GetNumberWaiting() int {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.round.count
}

func (b *cyclicBarrier) IsBroken() bool {
	b.lock.RLock()
	defer b.lock.RUnlock()

	return b.round.isBroken
}
