package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const (
	philosopherCount = 2
)

type Chopstick struct {
	sync.Mutex
}

func (c *Chopstick) take() {
	c.Lock()
}

func (c *Chopstick) leave() {
	c.Unlock()
}

func philospher(chopsticks []*Chopstick, index int, barrier *sync.WaitGroup) {
	if index == 0 {
		chopsticks[index+1].take()
		chopsticks[index].take()
	} else {
		chopsticks[index].take()
		chopsticks[(index+1)%len(chopsticks)].take()
	}
	fmt.Printf("El filósofo %d está comiendo\n", index)
	<-time.After(time.Millisecond * 300)
	fmt.Printf("El filósofo %d se va a pensar\n", index)
	chopsticks[index].leave()
	chopsticks[(index+1)%len(chopsticks)].leave()

	barrier.Done()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	barrier := sync.WaitGroup{}
	var chopsticks [philosopherCount + 1]*Chopstick

	for i := 0; i < philosopherCount+1; i++ {
		chopsticks[i] = &Chopstick{}
	}

	for i := 0; i < philosopherCount; i++ {
		i := i
		barrier.Add(1)
		go philospher(chopsticks[:], i, &barrier)
	}

	barrier.Wait()
}
