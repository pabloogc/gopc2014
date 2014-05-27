package app

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Client interface {
	Inbox() chan<- string
	Done() <-chan error
	Name() string
	Connect()
	Disconnect()
}

type clientImpl struct {
	sleepTime  int
	name       string
	inbox      chan (string)
	done       chan (error)
	disconnect chan (int)
	ticker     *time.Ticker
	friend     Client
	server     Server
}

func NewClient(name string, sleep int, server Server) Client {
	return &clientImpl{
		name:       name,
		server:     server,
		inbox:      make(chan string, 10),
		done:       make(chan error),
		disconnect: make(chan int, 1),
		ticker:     time.NewTicker(time.Millisecond * time.Duration(rand.Int()%sleep)),
		sleepTime:  sleep,
	}
}

func (c *clientImpl) Connect() {
	fmt.Printf("%s connecting \n", c.name)
	c.friend = c.server.FindFriend(c)
	fmt.Printf("%s paired with: %s \n", c.name, c.friend.Name())
	doneReading := make(chan int)
	doneWriting := make(chan int)
	go c.readForSomeTime(doneReading)
	go c.writeForSomeTime(doneWriting)
	<-doneReading
	<-doneWriting
	c.done <- nil //any error? send it here
}

func (c *clientImpl) readForSomeTime(done chan<- int) {
	defer func() {
		done <- 1
	}()

	timeOut := time.After(time.Second * 5)
	for {
		select {
		case <-timeOut:
			fmt.Printf("%s: Bored of reading. I'm leaving\n", c.name)
			c.friend.Disconnect()
			return
		case message, open := <-c.inbox:
			if !open {
				fmt.Printf("%s says goodbye!\n", c.friend.Name())
				return
			} else {
				fmt.Printf("%s: new message from %s -- %s\n", c.name, c.friend.Name(), message)
			}
		}
	}
}

func (c *clientImpl) writeForSomeTime(done chan<- int) {
	defer func() {
		close(c.friend.Inbox())
		done <- 1
	}()

	count := 0
	for {
		select {
		case <-c.disconnect:
			fmt.Printf("%s: %s left, nobody likes me u_u\n", c.name, c.friend.Name())
			return
		case <-c.ticker.C:
			count = count + 1
			select {
			case c.friend.Inbox() <- fmt.Sprintf("%d %s", count, c.name):
				if count > 10 {
					return
				}
			case <-c.disconnect:
				return
			}
		}
	}
}

func (c *clientImpl) Inbox() chan<- string {
	return c.inbox
}

func (c *clientImpl) Name() string {
	return c.name
}

func (c *clientImpl) Done() <-chan error {
	return c.done
}

func (c *clientImpl) Disconnect() {
	c.disconnect <- 1
	c.ticker.Stop()
}
