package cache

func (u *util[V]) Delete(k string) bool {
	u.lock.Lock()
	if _, found := u.data[k]; !found {
		u.lock.Unlock()
		return false
	}

	delete(u.data, k)
	ki := 0
	for i := range u.keyList {
		if u.keyList[i] == k {
			ki = i
			break
		}
	}
	u.keyList = append(u.keyList[:ki], u.keyList[ki+1:]...)

	u.lock.Unlock()
	return true
}
