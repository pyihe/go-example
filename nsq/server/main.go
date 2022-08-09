package main

import (
	"fmt"
	"time"

	nsqs "github.com/pyihe/go-example/nsq"
)

func main() {
	go func() {
		config := nsqs.NewConfig()
		config.DialTimeout = 10 * time.Second
		config.LookupdPollTimeout = 5 * time.Second
		config.ReadTimeout = 20 * time.Second
		config.WriteTimeout = 2 * time.Second
		config.HeartbeatInterval = 15 * time.Second
		producer, err := nsqs.AddProducer("192.168.1.192:4150", config)
		if err != nil {
			fmt.Printf("add producer err1: %v\n", err)
			return
		}
		defer producer.Stop()
		for {
			if err = producer.Publish("topic.game", []byte("topic.game")); err != nil {
				fmt.Printf("publish err1: %v\n", err)
				return
			}
			fmt.Printf("topic.game send\n")
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		config := nsqs.NewConfig()
		config.DialTimeout = 10 * time.Second
		config.LookupdPollTimeout = 5 * time.Second
		config.ReadTimeout = 20 * time.Second
		config.WriteTimeout = 2 * time.Second
		config.HeartbeatInterval = 15 * time.Second
		producer, err := nsqs.AddProducer("192.168.1.77:4150", config)
		if err != nil {
			fmt.Printf("add producer err2: %v\n", err)
			return
		}
		defer producer.Stop()
		for {
			if err = producer.Publish("topic.user", []byte("topic.user")); err != nil {
				fmt.Printf("publish err2: %v\n", err)
				return
			}
			fmt.Printf("topic.user send\n")
			time.Sleep(1 * time.Second)
		}
	}()

	select {}
}
