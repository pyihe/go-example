package etcd

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	etcdClient *clientv3.Client
}

func New(config clientv3.Config) *Client {
	c, err := clientv3.New(config)
	if err != nil {
		panic(err)
	}
	return &Client{
		etcdClient: c,
	}
}

func (c *Client) Close() error {
	return c.etcdClient.Close()
}

func (c *Client) Register(k, v string) error {
	var kv = clientv3.NewKV(c.etcdClient)
	var lease = clientv3.NewLease(c.etcdClient)
	var leaseRsp, err = lease.Grant(context.Background(), 10)
	if err != nil {
		return err
	}

	_, err = kv.Put(context.Background(), k, v, clientv3.WithLease(leaseRsp.ID))
	if err != nil {
		return err
	}

	// put并且绑定lease后，需要KeepAlive API来保持服务的正常状态
	aliveChan, err := lease.KeepAlive(context.Background(), leaseRsp.ID)
	if err != nil {
		return err
	}

	go func(ch <-chan *clientv3.LeaseKeepAliveResponse) {
		for data := range ch {
			_ = data
		}
	}(aliveChan)

	return nil
}

func (c *Client) Watch(key string, handler func(event *clientv3.Event), opts ...clientv3.OpOption) {
	go func() {
		var watcher = clientv3.NewWatcher(c.etcdClient)
		defer watcher.Close()

		var watchChan = watcher.Watch(context.Background(), key, opts...)
		for data := range watchChan {
			if data.Canceled {
				break
			}
			for _, event := range data.Events {
				if handler != nil {
					handler(event)
				}
			}
		}
	}()
}

func (c *Client) UnRegister(key string, opts ...clientv3.OpOption) error {
	var kv = clientv3.NewKV(c.etcdClient)
	var _, err = kv.Delete(context.Background(), key, opts...)
	return err
}
