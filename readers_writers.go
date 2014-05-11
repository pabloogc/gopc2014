package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const (
	readerWriterCount = 1
)

type Access struct {
	readersWaiting int
	readersWorking int
	writersWaiting int
	writerWorking  bool
	exclusion      sync.Locker
	readerAccess   sync.Cond
	writerAccess   sync.Cond
}

func (a *Access) init() {
	a.readersWaiting = 0
	a.readersWorking = 0
	a.writersWaiting = 0
	a.writerWorking = false
	a.exclusion = new(sync.Mutex)
	a.readerAccess = *sync.NewCond(a.exclusion)
	a.writerAccess = *sync.NewCond(a.exclusion)
}

func (a *Access) startReading() {
	a.exclusion.Lock()
	for a.writerWorking || a.writersWaiting > 0 {
		a.readersWaiting++
		a.readerAccess.Wait()
		a.readersWaiting--
		if !a.writerWorking {
			break
		}
	}
	a.readersWorking++
	a.exclusion.Unlock()
}

func (a *Access) finishReading() {
	a.exclusion.Lock()
	a.readersWorking++
	if a.readersWorking == 0 && a.writersWaiting > 0 {
		a.writerAccess.Signal()
	}
	a.exclusion.Unlock()
}

func (a *Access) startWriting() {
	a.exclusion.Lock()
	for a.writerWorking || a.readersWorking > 0 {
		a.writersWaiting++
		a.writerAccess.Wait()
		a.writersWaiting--
	}
	a.writerWorking = true
	a.exclusion.Unlock()
}

func (a *Access) finishWriting() {
	a.exclusion.Lock()
	a.writerWorking = false
	if a.readersWaiting > 0 {
		a.readerAccess.Broadcast()
	} else if a.writersWaiting > 0 {
		a.writerAccess.Signal()
	}
	a.exclusion.Unlock()
}

func reader(access *Access, index int, barrier *sync.WaitGroup) {
	access.startReading()
	//Leer dato
	fmt.Printf("El lector %d está leyendo\n", index)
	access.finishReading()
	//Procesar dato
	<-time.After(time.Millisecond * 300)
	fmt.Printf("El lector %d está procesando el dato\n", index)

	barrier.Done()
}

func writer(access *Access, index int, barrier *sync.WaitGroup) {
	// Generar dato
	fmt.Printf("El escritor %d está generando el dato\n", index)
	<-time.After(time.Millisecond * 300)
	access.startWriting()
	// Escribir dato
	fmt.Printf("El escritor %d está escribiendo el dato\n", index)
	access.finishWriting()

	barrier.Done()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	barrier := sync.WaitGroup{}

	access := &Access{}
	access.init()

	for i := 0; i < readerWriterCount; i++ {
		i := i
		barrier.Add(2)
		go reader(access, i, &barrier)
		go writer(access, i, &barrier)
	}

	barrier.Wait()
}
