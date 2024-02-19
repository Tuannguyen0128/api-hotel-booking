package memorizer

import (
	"context"
	"sync"
	"time"
)

type Util[V interface{}] interface {
	Get(ctx context.Context, k string) (V, error)
	Delete(k string)
}

type util[V interface{}] struct {
	cache   cache[V]
	getter  func(context.Context, string) (V, error)
	getting map[string]bool
	lock    sync.Mutex
	sleep   time.Duration
}

type cache[V interface{}] interface {
	Get(k string) (V, bool)
	Put(v V)
	Delete(k string) bool
}

func NewMemorizer[V interface{}](sleepInMS int, cache cache[V], getter func(context.Context, string) (V, error)) Util[V] {
	return &util[V]{
		cache:   cache,
		getter:  getter,
		getting: make(map[string]bool),
		lock:    sync.Mutex{},
		sleep:   time.Millisecond * time.Duration(sleepInMS),
	}
}
