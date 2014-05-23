package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	count := 30
	sem := NewSemaphore(1)
	wait := NewSemaphore(0)

	for i := 0; i < count; i++ {
		go thread(sem, wait, i)
	}

	for i := 0; i < count; i++ {
		wait.Acquire()
	}
}

func thread(sem, wait *Semaphore, index int) {
	<-time.After(time.Duration(rand.Int() % 10000))
	sem.Acquire()
	fmt.Printf("Hilo %d entra en la rc\n", index)
	<-time.After(time.Millisecond * 30)
	fmt.Printf("Hilo %d sale  en la rc\n", index)
	fmt.Println()
	sem.Release()
	wait.Release()
}

type Semaphore struct {
	lock chan int
}

func NewSemaphore(max int) *Semaphore {
	sem := &Semaphore{lock: make(chan int, max)}
	for i := 0; i < max; i++ {
		sem.lock <- 1
	}
	return sem
}

func (s *Semaphore) Acquire() {
	<-s.lock
}

func (s *Semaphore) Release() {
	s.lock <- 0
}
