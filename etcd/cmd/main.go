package main

import (
	"fmt"

	"github.com/pyihe/go-example/etcd"
	"github.com/pyihe/go-pkg/tools"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	client := etcd.New(clientv3.Config{Endpoints: []string{"192.168.1.77:2379"}})
	defer client.Close()

	err := client.Register("test/put", "hahaha")
	if err != nil {
		fmt.Printf("register err: %v\n", err)
		return
	}

	client.Watch("test", handler, clientv3.WithPrefix())

	tools.Wait()
}

func handler(event *clientv3.Event) {
	fmt.Println(event.IsModify(), event.IsCreate())
	fmt.Printf("kv发生变化: <%v,%v>\n", string(event.Kv.Key), string(event.Kv.Value))
}
