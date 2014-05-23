package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	readerWriterCount = 10
)

type Access struct {
	readersWaiting int
	readersWorking int
	writersWaiting int
	writerWorking  bool
	exclusion      sync.Locker
	readerAccess   *sync.Cond
	writerAccess   *sync.Cond
}

func NewAccess() *Access {
	a := new(Access)
	a.readersWaiting = 0
	a.readersWorking = 0
	a.writersWaiting = 0
	a.writerWorking = false
	a.exclusion = new(sync.Mutex)
	a.readerAccess = sync.NewCond(a.exclusion)
	a.writerAccess = sync.NewCond(a.exclusion)
	return a
}

func (a *Access) startReading() {
	a.exclusion.Lock()
	if a.writerWorking || a.writersWaiting > 0 {
		a.readersWaiting++
		a.readerAccess.Wait()
		a.readersWaiting--
	}
	a.readersWorking++
	a.exclusion.Unlock()
}

func (a *Access) finishReading() {
	a.exclusion.Lock()
	a.readersWorking--
	if a.readersWorking == 0 && a.writersWaiting > 0 {
		a.writerAccess.Signal()
	}
	a.exclusion.Unlock()
}

func (a *Access) startWriting() {
	a.exclusion.Lock()
	if a.writerWorking || a.readersWorking > 0 {
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
	<-time.After(time.Millisecond * time.Duration(rand.Int()%1000))
	access.startReading()

	//Leer dato
	fmt.Printf("El lector %d está leyendo\n", index)
	<-time.After(time.Millisecond * 200)
	access.finishReading()

	//Procesar dato
	fmt.Printf("El lector %d está procesando el dato\n", index)

	barrier.Done()
}

func writer(access *Access, index int, barrier *sync.WaitGroup) {
	// Generar dato
	<-time.After(time.Millisecond * 10)
	fmt.Printf("El escritor %d ha generado el dato\n", index)

	access.startWriting()

	// Escribir dato
	fmt.Printf("El escritor %d está escribiendo el dato\n", index)
	<-time.After(time.Millisecond * 50)

	access.finishWriting()

	barrier.Done()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	barrier := sync.WaitGroup{}

	access := NewAccess()

	for i := 0; i < readerWriterCount; i++ {
		barrier.Add(2)
		go writer(access, i, &barrier)
		go reader(access, i, &barrier)

	}
	barrier.Wait()
}
