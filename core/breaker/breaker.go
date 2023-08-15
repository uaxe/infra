package breaker

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

/*
	熔断器
    Google SRE 弹性熔断算法，弹性熔断是根据成功率动态调整的，当成功率越高的时候，
	被熔断的概率就越小；反之，当成功率越低时，被熔断的概率就相应增大。
	基于google的sre算法实现了熔断器逻辑，并在redis等客户端操作的时候引入了熔断器
	参考：
    熔断原理与实现Golang版 https://blog.csdn.net/jfwan/article/details/109328874
	Google SRE 弹性熔断算法实现分析:    https://pandaychen.github.io/2020/05/10/A-GOOGLE-SRE-BREAKER/
	代码参考： https://github.com/tal-tech/go-zero/blob/master/core/breaker/breaker.go
*/

const (
	numHistoryReasons = 5
	timeFormat        = "15:04:05"
)

var ErrServiceUnavailable = errors.New("circuit breaker is open")

type (
	// 检查错误是否符合预期
	Acceptable func(err error) bool

	// 熔断器
	Breaker interface {
		// 返回熔断器名称
		Name() string

		// 检查请求是否被允许
		Allow() (Promise, error)

		Do(req func() error) error

		DoWithAcceptable(req func() error, acceptable Acceptable) error

		DoWithFallback(req func() error, fallback func(err error) error) error

		DoWithFallbackAcceptable(req func() error, fallback func(err error) error, acceptable Acceptable) error
	}

	Option func(breaker *circuitBreaker)

	Promise interface {
		Accept()
		Reject(reason string)
	}

	internalPromise interface {
		Accept()
		Reject()
	}

	circuitBreaker struct {
		name string
		throttle
	}

	internalThrottle interface {
		allow() (internalPromise, error)
		doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error
	}

	throttle interface {
		allow() (Promise, error)
		doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error
	}
)

func NewBreaker(opts ...Option) Breaker {
	var b circuitBreaker
	for _, opt := range opts {
		opt(&b)
	}
	if len(b.name) == 0 {
		b.name = time.Now().Format(timeFormat)
	}
	b.throttle = newLoggedThrottle(b.name, newGoogleBreaker())
	return &b
}

func defaultAcceptable(err error) bool {
	return err == nil
}

func (cb *circuitBreaker) Allow() (Promise, error) {
	return cb.throttle.allow()
}

func (cb *circuitBreaker) Do(req func() error) error {
	return cb.throttle.doReq(req, nil, defaultAcceptable)
}

func (cb *circuitBreaker) DoWithAcceptable(req func() error, acceptable Acceptable) error {
	return cb.throttle.doReq(req, nil, acceptable)
}

func (cb *circuitBreaker) DoWithFallback(req func() error, fallback func(err error) error) error {
	return cb.throttle.doReq(req, fallback, defaultAcceptable)
}

func (cb *circuitBreaker) DoWithFallbackAcceptable(req func() error, fallback func(err error) error,
	acceptable Acceptable) error {
	return cb.throttle.doReq(req, fallback, acceptable)
}

func (cb *circuitBreaker) Name() string {
	return cb.name
}

func WithName(name string) Option {
	return func(b *circuitBreaker) {
		b.name = name
	}
}

type loggedThrottle struct {
	name string
	internalThrottle
	errWin *errorWindow
}

func newLoggedThrottle(name string, t internalThrottle) loggedThrottle {
	return loggedThrottle{
		name:             name,
		internalThrottle: t,
		errWin:           new(errorWindow),
	}
}

type promiseWithReason struct {
	promise internalPromise
	errWin  *errorWindow
}

func (p promiseWithReason) Accept() {
	p.promise.Accept()
}

func (p promiseWithReason) Reject(reason string) {
	p.errWin.add(reason)
	p.promise.Reject()
}

func (lt loggedThrottle) allow() (Promise, error) {
	promise, err := lt.internalThrottle.allow()
	return promiseWithReason{
		promise: promise,
		errWin:  lt.errWin,
	}, lt.logError(err)
}

func (lt loggedThrottle) doReq(req func() error, fallback func(err error) error, acceptable Acceptable) error {
	return lt.logError(lt.internalThrottle.doReq(req, fallback, func(err error) bool {
		accept := acceptable(err)
		if !accept {
			lt.errWin.add(err.Error())
		}
		return accept
	}))
}

func (lt loggedThrottle) logError(err error) error {
	if err == ErrServiceUnavailable {
		log.Printf("stdout", "%s breaker is open and requests dropped\nlast errors:\n %s", lt.name, lt.errWin)
	}
	return err
}

type errorWindow struct {
	reasons [numHistoryReasons]string
	index   int
	count   int
	lock    sync.Mutex
}

func (ew *errorWindow) add(reason string) {
	ew.lock.Lock()
	ew.reasons[ew.index] = fmt.Sprintf("%s %s", Time().Format(timeFormat), reason)
	ew.index = (ew.index + 1) % numHistoryReasons
	ew.count = MinInt(ew.count+1, numHistoryReasons)
	ew.lock.Unlock()
}

func (ew *errorWindow) String() string {
	var reasons []string

	ew.lock.Lock()
	// reverse order
	for i := ew.index - 1; i >= ew.index-ew.count; i-- {
		reasons = append(reasons, ew.reasons[(i+numHistoryReasons)%numHistoryReasons])
	}
	ew.lock.Unlock()

	return strings.Join(reasons, "\n")
}
