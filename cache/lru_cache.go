package cache

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/uaxe/infra/schedule"
)

type Entry struct {
	Timerid uint32
	Value   any
}

type ChanMessage struct {
	Key any `json:"k"`
}

type LRUCache struct {
	ctx   context.Context
	hit   uint64
	total uint64
	cache *cache
	tw    *schedule.TimerWheel
}

func NewLRUCache(
	ctx context.Context,
	maxcapacity int,
	expiredTw *schedule.TimerWheel,
	OnEvicted func(k, v any)) *LRUCache {

	c := New(maxcapacity)
	c.OnEvicted = func(key, value any) {
		vv := value.(Entry)
		if nil != OnEvicted {
			OnEvicted(key.(any), vv.Value)
		}
		expiredTw.CancelTimer(vv.Timerid)
	}

	lru := &LRUCache{
		ctx:   ctx,
		cache: c,
		tw:    expiredTw}
	return lru
}

func (l *LRUCache) HitRate() (int, int) {
	currHit := l.hit
	currTotal := l.total

	if currTotal <= 0 {
		return 0, l.Length()
	}
	return int(currHit * 100 / currTotal), l.Length()
}

func (l *LRUCache) Get(key any) (any, bool) {
	atomic.AddUint64(&l.total, 1)
	if v, ok := l.cache.Get(key); ok {
		atomic.AddUint64(&l.hit, 1)
		return v.(Entry).Value, true
	}
	return nil, false
}

func (l *LRUCache) onMessage(key any) {
	l.Remove(key)
}

func (l *LRUCache) Put(key, v any, ttl time.Duration) chan time.Time {
	vv := Entry{Value: v}
	var ttlChan chan time.Time
	if ttl > 0 {
		if nil != l.tw {
			if val, ok := l.cache.Get(key); ok {
				if exist, ok := val.(Entry); ok {
					l.tw.CancelTimer(exist.Timerid)
				}
			}
			timerid, ch := l.tw.AddTimer(ttl, func(t time.Time) {
				l.cache.Remove(key)
			}, func(t time.Time) {

			})
			vv.Timerid = timerid
			ttlChan = ch
		}
	}
	l.cache.Add(key, vv)
	return ttlChan
}

func (l *LRUCache) Remove(key any) any {
	vv := l.cache.Remove(key)
	if vv != nil {
		e := vv.(Entry)
		if nil != l.tw {
			l.tw.CancelTimer(e.Timerid)
		}
		return e.Value
	}

	return nil
}

func (l *LRUCache) Contains(key any) bool {
	if _, ok := l.cache.Get(key); ok {
		return ok
	}
	return false
}

func (l *LRUCache) Length() int {
	return l.cache.Len()
}

func (l *LRUCache) Iterator(do func(k, v any) error) {
	l.cache.Iterator(func(k, v any) error {
		vv := v.(Entry)
		return do(k, vv.Value)
	})
}
