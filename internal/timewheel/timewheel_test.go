package timewheel

import (
	"go-cache/log"
	"testing"
	"time"
)

func TestNewTimeWheel(t *testing.T) {
	tw := NewTimeWheel(time.Second, 5)
	tw.Start()

	time.Sleep(6)

	log.Logger.Info("over")

	tw.AddJob("1", time.Second*2, func() {
		log.Logger.Info("1")
	})

	tw.AddJob("2", time.Second, func() {
		log.Logger.Info("2")
	})

	tw.AddJob("3", time.Second*3, func() {
		log.Logger.Info("3")
	})

	tw.AddJob("8", time.Second*8, func() {
		log.Logger.Info("8")
	})
	tw.AddJob("9", time.Second*9, func() {
		log.Logger.Info("9")
	})
	// 防止指令重排
	time.Sleep(1 * time.Second)
	tw.RemoveJob("8")

	time.Sleep(20 * time.Second)
}
