package rand

import (
	"sync"
	"testing"
)

func TestRandConcurrent(t *testing.T) {
	const N = 1000
	ch := make(chan int64, N)
	seen := make(map[int64]struct{}, N)
	empty := struct{}{}
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			ch <- Rand.Int63()
		}()
	}
	wg.Wait()
	close(ch)
	for v := range ch {
		seen[v] = empty
	}
	if len(seen) != N {
		t.Errorf("got %d different numbers, want %d", len(seen), N)
	}
}
