package main

import (
	"github.com/pabloogc/gopc2014/app"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	server := app.NewServer()
	client1 := app.NewClient("Pablo", 100, server)
	client2 := app.NewClient("Izan", 1000, server)
	go client1.Connect()
	go client2.Connect()

	<-client1.Done()
	<-client2.Done()

}
