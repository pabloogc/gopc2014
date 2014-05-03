package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

//*********************

func main() {
	//use all the cores
	runtime.GOMAXPROCS(runtime.NumCPU())

	//random seed
	rand.Seed(time.Now().UnixNano())

	end := make(chan bool)
	c := make(chan string)
	go consumer(c, end)
	go producer(c, "c1")
	go producer(c, "c2")
	//panic("hola")
	<-end
}

func producer(c chan<- string, nombre string) {
	for i := 0; i < 10; i++ {
		c <- fmt.Sprintf("%s %d", nombre, i)
		<-time.After(time.Millisecond * 100)
	}
}

func consumer(c <-chan string, f chan<- bool) {
	timeOut := time.After(time.Millisecond * 300)
	for {
		select {
		case <-timeOut:
			fmt.Println("timeout")
			f <- true
			return
		case s := <-c:
			fmt.Println(s)
		}
	}
}
