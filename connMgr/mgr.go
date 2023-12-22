package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"
)

var (
	ErrConnectionUnavailable = errors.New("connection unavailable")
)

type Dialer func(network, address string) (net.Conn, error)

type AfterFunc func(duration time.Duration) <-chan time.Time

type Manager struct {
	dialer  Dialer        // “拨号器”
	network string        // 连接类型
	address string        // 连接地址
	after   AfterFunc     // 重连func
	takec   chan net.Conn // 用于获取连接
	putc    chan error    // 接收来自连接的错误，如果出错，则重新连接
}

func NewManager(d Dialer, network, address string, after AfterFunc) *Manager {
	m := &Manager{
		dialer:  d,
		network: network,
		address: address,
		after:   after,
		takec:   make(chan net.Conn),
		putc:    make(chan error),
	}

	go m.loop()

	return m
}

func NewDefaultManager(network, address string) *Manager {
	return NewManager(net.Dial, network, address, time.After)
}

func (m *Manager) Take() net.Conn {
	return <-m.takec
}

func (m *Manager) Put(err error) {
	m.putc <- err
}

func (m *Manager) Write(b []byte) (int, error) {
	conn := m.Take()
	if conn == nil {
		return 0, ErrConnectionUnavailable
	}
	n, err := conn.Write(b)
	defer m.Put(err)
	return n, err
}

func (m *Manager) loop() {
	var (
		conn          = dial(m.dialer, m.network, m.address) // 新建连接
		connChan      = make(chan net.Conn, 1)               // 传递连接的Channel
		reconnectChan <-chan time.Time                       // 重连的Channel
		backoff       = time.Second                          // 重连的时间间隔
	)

	// 将新建的连接传递进connChan
	connChan <- conn

	for {
		select {
		case <-reconnectChan:
			fmt.Println("重连...")
			reconnectChan = nil
			go func() {
				connChan <- dial(m.dialer, m.network, m.address)
			}()
		case conn = <-connChan:
			if conn == nil {
				fmt.Println("conn为空, 需要重连...")
				backoff = exponential(backoff)
				reconnectChan = m.after(backoff)
			} else {
				fmt.Println("conn不为空，不需要重连...")
				backoff = time.Second
				reconnectChan = nil
			}
		case m.takec <- conn:
			fmt.Println("m.takec <- conn")
		case err := <-m.putc:
			fmt.Println("收到err: ", err)
			if err != nil && conn != nil {
				if conn != nil {
					conn.Close()
				}
				conn = nil
				reconnectChan = m.after(time.Nanosecond)
			}
		}
	}
}

func dial(d Dialer, network, address string) net.Conn {
	conn, err := d(network, address)
	if err != nil {
		conn = nil
	}
	return conn
}

func exponential(d time.Duration) time.Duration {
	d *= 2
	jitter := rand.Float64() + 0.5
	d = time.Duration(int64(float64(d.Nanoseconds()) * jitter))
	if d > time.Minute {
		d = time.Minute
	}
	return d
}
