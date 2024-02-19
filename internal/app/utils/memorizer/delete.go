package memorizer

func (u *util[V]) Delete(k string) {
	u.lock.Lock()
	_ = u.cache.Delete(k)
	u.lock.Unlock()
}
