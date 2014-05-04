package app

import (
	"sync"
)

type Server interface {
	findFriend(Client) Client
}
type serverImpl struct {
	turn  bool
	c1    chan (Client)
	c2    chan (Client)
	mutex sync.Locker
}

func NewServer() Server {
	return &serverImpl{
		turn:  false,
		c1:    make(chan Client, 1),
		c2:    make(chan Client, 1),
		mutex: new(sync.Mutex),
	}
}

func (s *serverImpl) findFriend(client Client) Client {
	s.turn = !s.turn //sneaky!
	if s.turn {
		s.c1 <- client
		return <-s.c2
	} else {
		s.c2 <- client
		return <-s.c1
	}
}
