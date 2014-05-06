package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const (
	count = 3
)

var (
	finAC = make(chan int, count-1)
	finBC [count]chan string
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	wait := new(sync.WaitGroup)

	for i := 0; i < count; i++ {
		i := i
		wait.Add(1)
		finBC[i] = make(chan string)
		go thread(i, wait)
	}

	wait.Wait()
	fmt.Println()
}

func thread(index int, wait *sync.WaitGroup) {
	for i := 0; i < 3; i++ {
		fmt.Printf("%dA ", index)
		<-time.After(time.Millisecond * 300)
		finA(index)
		fmt.Printf("%dB ", index)
		<-time.After(time.Millisecond * 300)
		finB(index)
		fmt.Printf("%dC ", index)
		<-time.After(time.Millisecond * 300)
		finC(index, i)
	}
	wait.Done()
}

func finA(index int) {
	if index == 0 {
		for i := 0; i < count-1; i++ {
			<-finAC //acquire
		}
	} else {
		finAC <- 1 //release
	}
}

func finB(index int) {
	if index == 0 {
		for i := 1; i < count; i++ {
			finBC[i] <- "fin!" //acquire
		}
	} else {
		<-finBC[index] //release
	}
}

func finC(index int, id int) {
}
