package cache

import "time"

func (u *util[V]) Put(v V) {
	k := v.GetId()
	u.lock.Lock()
	if _, existed := u.data[k]; !existed {
		u.keyList = append(u.keyList, k)
	}
	u.data[k] = item[V]{v: v, timeout: time.Now().Add(u.timeout)}

	// ensure size
	for len(u.keyList) > u.size {
		delete(u.data, u.keyList[0])
		u.keyList = u.keyList[1:]
	}
	u.lock.Unlock()
}
