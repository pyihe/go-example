package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
)

var (
	ErrEmptyServerList = errors.New("server is empty")
)

var cacheNodes = []string{
	"127.0.0.1:12000",
	"127.0.0.1:12001",
	"127.0.0.1:12003",
}

type roundSelector struct {
	mu      sync.Mutex
	servers []net.Addr
	idx     int
}

func newRoundSelector(addrs ...string) (memcache.ServerSelector, error) {
	n := len(addrs)
	if n == 0 {
		return nil, ErrEmptyServerList
	}
	naddr := make([]net.Addr, len(addrs))
	for i, addr := range addrs {
		var sAddr net.Addr
		if strings.Contains(addr, "/") {
			unixAddr, err := net.ResolveUnixAddr("unix", addr)
			if err != nil {
				return nil, err
			}
			sAddr = unixAddr
		} else {
			tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
			if err != nil {
				return nil, err
			}
			sAddr = tcpAddr
		}
		naddr[i] = sAddr
	}

	return &roundSelector{
		mu:      sync.Mutex{},
		servers: naddr,
		idx:     0,
	}, nil
}

func (rs *roundSelector) PickServer(key string) (net.Addr, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	n := len(rs.servers)
	if n == 0 {
		return nil, ErrEmptyServerList
	}
	addr := rs.servers[rs.idx%n]
	rs.idx += 1
	if rs.idx == n {
		rs.idx = 0
	}
	return addr, nil
}

func (rs *roundSelector) Each(fn func(addr net.Addr) error) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	for _, addr := range rs.servers {
		if err := fn(addr); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	client := memcache.New(cacheNodes...)
	defer client.Close()

	items := []*memcache.Item{
		{
			Key:        "k1",
			Value:      []byte{1},
			Expiration: 60,
		},
		{
			Key:   "k2",
			Value: []byte{2},
		},
	}

	if err := client.Set(items[0]); err != nil {
		fmt.Printf("set err: %v\n", err)
		return
	}
	if err := client.Add(items[1]); err != nil {
		fmt.Printf("add err: %v\n", err)
		return
	}
	if item, err := client.Get(items[0].Key); err != nil {
		fmt.Printf("get err: %v\n", err)
		return
	} else {
		fmt.Printf("get item: %+v\n", item)
	}

	if err := client.Delete(items[0].Key); err != nil {
		fmt.Printf("del err: %v\n", err)
		return
	}

	// 自定义selector
	//selector, err := newRoundSelector(cacheNodes...)
	//if err != nil {
	//	fmt.Printf("new selector err: %v\n", err)
	//	return
	//}
	//client = memcache.NewFromSelector(selector)
}
