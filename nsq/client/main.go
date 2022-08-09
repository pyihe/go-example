package main

import (
	"fmt"

	"github.com/nsqio/go-nsq"
	nsqs "github.com/pyihe/go-example/nsq"
	"github.com/pyihe/go-pkg/bytes"
)

func main() {
	nsqs.NewConsumerMgr(nsqs.NewConfig(), []string{"192.168.1.77:4161", "192.168.1.192:4161"})

	if err := nsqs.AddConsumer("topic.game", "channel1", &gameC1{}); err != nil {
		fmt.Printf("add consumer err1: %v\n", err)
		return
	}
	if err := nsqs.AddConsumer("topic.game", "channel2", &gameC2{}); err != nil {
		fmt.Printf("add consumer err2: %v\n", err)
		return
	}

	if err := nsqs.AddConsumer("topic.user", "channel1", &userC1{}); err != nil {
		fmt.Printf("errxx: %v\n", err)
		return
	}
	if err := nsqs.AddConsumer("topic.user", "channel2", &userC2{}); err != nil {
		fmt.Printf("errxx: %v\n", err)
		return
	}
	fmt.Printf("添加成功...\n")
	select {}
}

type gameC1 struct{}

func (g *gameC1) HandleMessage(message *nsq.Message) error {
	fmt.Printf("topic.game channel1: %v\n", bytes.String(message.Body))
	return nil
}

type gameC2 struct{}

func (g *gameC2) HandleMessage(message *nsq.Message) error {
	fmt.Printf("topic.game channel2: %v\n", bytes.String(message.Body))
	return nil
}

type userC1 struct{}

func (g *userC1) HandleMessage(message *nsq.Message) error {
	fmt.Printf("topic.user channel1: %v\n", bytes.String(message.Body))
	return nil
}

type userC2 struct{}

func (g *userC2) HandleMessage(message *nsq.Message) error {
	fmt.Printf("topic.user channel2: %v\n", bytes.String(message.Body))
	return nil
}
