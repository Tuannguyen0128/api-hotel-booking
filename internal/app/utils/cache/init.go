package cache

import (
	"sync"
	"time"
)

// TODO Put existed key do not shift key FIFO list
type getId interface {
	GetId() string
}

type Util[V getId] interface {
	Delete(k string) bool
	Get(k string) (V, bool)
	Put(v V)

	inspect() string
}

type util[V getId] struct {
	timeout time.Duration
	data    map[string]item[V]
	keyList []string
	size    int
	lock    sync.Mutex

	hit  int
	miss int
}

type item[V getId] struct {
	v       V
	timeout time.Time
}

func NewCache[V getId](size, timeoutInSec int) Util[V] {
	return &util[V]{
		timeout: time.Second * time.Duration(timeoutInSec),
		data:    make(map[string]item[V]),
		keyList: []string{},
		size:    size,
		lock:    sync.Mutex{},
	}
}
