package cache

import "time"

func (u *util[V]) Get(k string) (V, bool) {
	var v V
	u.lock.Lock()
	i, found := u.data[k]
	u.lock.Unlock()
	if found {
		if time.Now().Before(i.timeout) {
			u.hit++
			return i.v, true
		}

		u.Delete(k)
	}

	u.miss++
	return v, false
}
