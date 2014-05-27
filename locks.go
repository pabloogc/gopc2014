package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var (
	lock      = new(sync.Mutex)
	lock2     = new(sync.Mutex)
	condition = sync.NewCond(lock)
	count     = 0
	num       = 2
	stop      = 1
	k         = 1
	barrier   = new(sync.WaitGroup)
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	for i := 0; i < num; i++ {
		barrier.Add(1)
		go func(i int) {
			name := fmt.Sprintf("h%d", i)
			nested(name)
			barrier.Done()
		}(i)
	}

	barrier.Wait()
}

func critic(name string) {
	lock.Lock()

	count++

	if count == stop {
		fmt.Println("    Esperando  -> 	" + name)
		condition.Wait()
		fmt.Println("    !!! Reanudado de forma inmediata -> " + name)
		if count != stop+k {
			panic("Imposible, alguien se ha colado!")
		}
	}

	time.Sleep(time.Duration(time.Millisecond * 100))

	fmt.Printf("%d/%d, completado %s\n", count, num, name)

	if count == stop+k {
		condition.Signal()
	}
	lock.Unlock()
}

func nested(name string) {
	lock2.Lock()
	<-time.After(time.Millisecond * 30)
	fmt.Printf("outer %s\n", name)
	//Deadlock! outer lock is never released when
	//blocking inside
	critic(name)
	lock2.Unlock()
}
