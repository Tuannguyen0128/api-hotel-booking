package memorizer

import (
	"context"
	"time"
)

func (u *util[V]) Get(ctx context.Context, k string) (V, error) {
	if v, found := u.cache.Get(k); found {
		return v, nil
	}

	if _, found := u.getting[k]; found {
		time.Sleep(u.sleep)
		return u.Get(ctx, k)
	}

	return u.doGet(ctx, k)
}

func (u *util[V]) doGet(ctx context.Context, k string) (V, error) {
	u.lock.Lock()
	if v, found := u.cache.Get(k); found {
		u.lock.Unlock()
		return v, nil
	}

	if _, found := u.getting[k]; found {
		u.lock.Unlock()
		time.Sleep(u.sleep)
		return u.Get(ctx, k)
	}

	u.getting[k] = true
	u.lock.Unlock()

	// real get
	v, err := u.getter(ctx, k)
	if err != nil {
		u.lock.Lock()
		delete(u.getting, k)
		u.lock.Unlock()
		return v, err
	} else {
		u.cache.Put(v)
	}

	u.lock.Lock()
	delete(u.getting, k)
	u.lock.Unlock()

	return v, nil
}
