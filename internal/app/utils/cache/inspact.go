package cache

import "fmt"

func (u *util[V]) inspect() string {
	return fmt.Sprintf("hit rate %.2f %% hit %d miss %d Size %d : %+v", float32(u.hit)*100/float32(u.miss), u.hit, u.miss, len(u.data), u.data)
}
