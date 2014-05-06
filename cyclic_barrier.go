package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	count = 100
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	done := new(sync.WaitGroup)
	barrier := NewCyclicBarrer(count)

	for i := 0; i < count; i++ {
		done.Add(1)
		go thread(barrier, done)
	}

	done.Wait()
	fmt.Println()
}

func thread(barrier *CyclicBarrier, done *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		<-time.After(time.Millisecond * time.Duration(rand.Int()%1000))
		fmt.Printf("%d ", i)
		barrier.Wait()
	}
	done.Done()
}

type CyclicBarrier struct {
	concurrency int
	waiting     int
	lock        chan int
	waitC       chan int
}

func NewCyclicBarrer(concurrency int) *CyclicBarrier {
	return &CyclicBarrier{
		concurrency: concurrency,
		waiting:     0,
		lock:        make(chan int, 1),
		waitC:       make(chan int),
	}
}

func (c *CyclicBarrier) Wait() {
	c.lock <- 1
	c.waiting++
	if c.waiting == c.concurrency {
		c.waiting = 0
		for i := 1; i < c.concurrency; i++ {
			c.waitC <- 1
		}
		fmt.Println()
		<-c.lock
	} else {
		<-c.lock
		<-c.waitC
	}

}
