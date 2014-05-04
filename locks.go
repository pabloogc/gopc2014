package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var (
	lock      = new(sync.Mutex)
	condition = sync.NewCond(lock)
	count     = 0
	num       = 10
	stop      = 5
	barrier   = new(sync.WaitGroup)
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	for i := 0; i < num; i++ {
		barrier.Add(1)
		go func(i int) {
			critic(fmt.Sprintf("%d", i))
			barrier.Done()
		}(i)
	}

	barrier.Wait()
}

func critic(name string) {
	lock.Lock()

	count = count + 1

	if count == stop {
		fmt.Println("....Esperando -> 	" + name)
		condition.Wait()
		fmt.Println("....Despertado -> 	" + name)
		if count != stop+1 {
			panic("Alguien se ha colado!")
		}
	}

	time.Sleep(time.Duration(time.Millisecond * 100))
	fmt.Printf("%d/%d, completado %s\n", count, num, name)

	if count == stop+1 {
		fmt.Println("....Despertando")
		condition.Signal()
	}
	lock.Unlock()
}
