package cache

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestConcurrent(t *testing.T) {
	now := time.Now()
	cache := NewCache[dataWithId](10, 1)

	wait := sync.WaitGroup{}
	N := 100
	for i := 0; i < N; i++ {
		k := fmt.Sprintf("id%d", i)
		go func() {
			wait.Add(1)

			cache.Put(dataWithId{
				id: fmt.Sprintf(k),
			})
			cache.Get(fmt.Sprintf("id%d", rand.Intn(N)))
			cache.Get(k)
			cache.Delete(k)
			cache.Put(dataWithId{
				id: fmt.Sprintf(k),
			})

			wait.Done()
		}()
	}
	wait.Wait()
	fmt.Println(cache.inspect())
	fmt.Printf("total run time %s\n", time.Now().Sub(now))
}
