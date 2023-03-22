package messageQueue

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// benchmarkSimpleMQ_Push
// 100000 并发Push ，耗时：36ms
func BenchmarkSimpleMQ_Push(b *testing.B) {

	// 模拟MySQL的处理
	mq := NewSimpleMQ(500, func(msg string) {
		time.Sleep(10 * time.Millisecond)
		//fmt.Printf(msg + " ")
	})
	// monitor
	var done = make(chan struct{})
	go func() {
		for {
			fmt.Println()
			fmt.Println()
			if mq.Len() == 0 {
				break
			}
			fmt.Printf("\033[1;31mqueue Len(): %v, queue cap(): %v\033[0m\n", mq.Len(), mq.que.Cap())
			fmt.Printf("\033[1;31mmsgChan len(): %v\033[0m\n", len(mq.msgChan))
			time.Sleep(1 * time.Second)
		}
		done <- struct{}{}
	}()

	// producer
	strat := time.Now()
	var wg sync.WaitGroup
	maxMsg := 100_000
	goroutineNum := 1000
	wg.Add(goroutineNum)
	for i := 0; i < goroutineNum; i++ {
		go func(i int) {
			for i2 := 0; i2 < maxMsg/goroutineNum; i2++ {
				mq.Push(fmt.Sprintf("producer:%v-%v", i, i2))
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	// all msg pushed
	spent := time.Since(strat)
	fmt.Printf("\033[1;32mall msg pushed, spent: %v, msgNum: %v\033[0m\n", spent, maxMsg)

	// wait for all msg consumed
	<-done
	fmt.Printf("\033[1;32mAll msg consumed\033[0m\n")
}
