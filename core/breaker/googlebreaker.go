package breaker

import (
	"math"
	"time"
)

const (
	// 250ms for bucket duration
	windowDuration = time.Second * 10
	buckets        = 40
	k              = 1.5
	protection     = 5
)

// googleBreaker is a netflixBreaker pattern from google.
// see Client-Side Throttling section in https://landing.google.com/sre/sre-book/chapters/handling-overload/
type googleBreaker struct {
	k     float64
	stat  *RollingWindow
	proba *Proba
}

func newGoogleBreaker() *googleBreaker {
	bucketDuration := time.Duration(int64(windowDuration) / int64(buckets))
	st := NewRollingWindow(buckets, bucketDuration)
	return &googleBreaker{
		stat:  st,
		k:     k,
		proba: NewProba(),
	}
}

func (b *googleBreaker) accept() error {
	// accepts为正常请求数，total为总请求数
	accepts, total := b.history()
	weightedAccepts := b.k * float64(accepts)
	// https://landing.google.com/sre/sre-book/chapters/handling-overload/#eq2101
	// 算法实现
	dropRatio := math.Max(0, (float64(total-protection)-weightedAccepts)/float64(total+1))
	if dropRatio <= 0 {
		return nil
	}
	// 是否超过比例
	if b.proba.TrueOnProba(dropRatio) {
		return ErrServiceUnavailable
	}
	return nil
}

func (b *googleBreaker) allow() (internalPromise, error) {
	if err := b.accept(); err != nil {
		return nil, err
	}
	return googlePromise{
		b: b,
	}, nil
}

// doReq 方法首先判断是否熔断，满足条件直接返回 error(circuit breaker is open)，不满足条件则对请求数进行累加
func (b *googleBreaker) doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error {
	if err := b.accept(); err != nil {
		if fallback != nil {
			return fallback(err)
		}
		return err
	}

	defer func() {
		if e := recover(); e != nil {
			b.markFailure()
			panic(e)
		}
	}()

	err := req()
	// 正常请求total和accepts都会加1
	if acceptable(err) {
		b.markSuccess()
	} else {
		// 请求失败只有total会加1
		b.markFailure()
	}
	return err
}

func (b *googleBreaker) history() (accepts int64, total int64) {
	b.stat.Reduce(func(b *Bucket) {
		accepts += int64(b.Sum)
		total += b.Count
	})
	return
}

func (b *googleBreaker) markSuccess() {
	b.stat.Add(1)
}

func (b *googleBreaker) markFailure() {
	b.stat.Add(0)
}

type googlePromise struct {
	b *googleBreaker
}

func (p googlePromise) Accept() {
	p.b.markSuccess()
}

func (p googlePromise) Reject() {
	p.b.markFailure()
}
