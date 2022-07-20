package main

import (
	"fmt"
	"net/rpc"
	"time"
)

func main() {
	s, err := serveTCP(":9999")
	if err != nil {
		panic(err)
	}

	game := new(Game)
	for i := 1; i <= 10; i++ {
		game.history.Store(int64(i), time.Now().Unix())
	}

	time.Sleep(1 * time.Second)

	if err = s.Register(game); err != nil {
		panic(err)
	}

	done := make(chan *rpc.Call, 10)
	go func() {
		login(done)
	}()
	for {
		select {
		case call := <-done:
			method := call.ServiceMethod
			if call.Error != nil {
				fmt.Printf("%s fail: %v\n", method, call.Error)
			}
			switch method {
			case "Game.Login":
				fmt.Println(*call.Reply.(*int64))
			}
		}
	}
}

func login(done chan *rpc.Call) {
	c, err := rpc.Dial("tcp", ":9999")
	if err != nil {
		panic(err)
	}
	for i := 1; i <= 10; i++ {
		go func(id int) {
			call := c.Go("Game.Login", id, new(int64), nil)
			done <- <-call.Done
		}(i)
	}
}
