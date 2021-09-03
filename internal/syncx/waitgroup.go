package syncx

import (
	"sync"
	"time"
)

type WaitGroup struct {
	wg sync.WaitGroup
}

func (w *WaitGroup) Add(delta int) {
	w.wg.Add(delta)
}

func (w *WaitGroup) Done() {
	w.wg.Done()
}

func (w *WaitGroup) Wait() {
	w.wg.Wait()
}
func (w *WaitGroup) WaitTime(duration time.Duration) bool {
	c := make(chan struct{}, 1)

	go func() {

		defer close(c)
		w.wg.Wait()
		c <- struct{}{}

	}()

	select {
	case <-c:
		return false
	case <-time.After(duration):
		return true
	}

}
