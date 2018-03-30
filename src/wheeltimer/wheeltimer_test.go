package wheeltimer

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"context"
)

func Test_NewTimer(t *testing.T) {
	wg := &sync.WaitGroup{}

	wheel := NewTimingWheel(context.TODO())
	timeout := NewOnTimeOut(func() {
		fmt.Println("需要执行的任务")
	})
	timerID := wheel.AddTimer(`1 * * * * *`, `unqiuename`, timeout)
	timerIDs := wheel.AddTimer(`2 * * * * *`, `caochun-task`, timeout)
	fmt.Printf("添加任务，任务ID：%d\n", timerID)
	fmt.Printf("添加任务，任务ID：%d\n", timerIDs)
	c := time.After(9 * time.Minute)
	wg.Add(1)
	go func() {
		for {
			select {
			case timeout := <-wheel.TimeOutChannel():
				timeout.Callback()
			case <-c:
				wg.Done()
			}
		}
	}()
	wg.Wait()
	wheel.Stop()
}
