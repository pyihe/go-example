package main

import (
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

func main() {
	var (
		tickChan = make(chan time.Time)
		after    = func(d time.Duration) <-chan time.Time {
			return time.After(1 * time.Millisecond)
			//return tickChan
		}
		dialConn = &mockConn{}
		dialErr  = error(nil)
		dialer   = func(string, string) (net.Conn, error) {
			return dialConn, dialErr
		}
		mgr = NewManager(dialer, "netw", "addr", after)
	)

	conn := mgr.Take()
	if conn == nil {
		fmt.Printf("nil conn\n")
		return
	}
	if _, err := conn.Write([]byte{1, 2, 3}); err != nil {
		fmt.Printf("write fail: %v\n", err)
		return
	}
	if want, have := uint64(3), atomic.LoadUint64(&dialConn.wr); want != have {
		fmt.Printf("want: %d, have: %d\n", want, have)
		return
	}
	mgr.Put(errors.New("should kill the connection"))

	for i := 0; i < 10; i++ {
		if conn = mgr.Take(); conn != nil {
			fmt.Printf("iteration %d: want nil connnection, but got real conn\n", i)
			return
		}
	}

	tickChan <- time.Now()
	if !within(100*time.Millisecond, func() bool {
		conn = mgr.Take()
		return conn != nil
	}) {
		fmt.Printf("conn remained nil\n")
		return
	}

	if _, err := conn.Write([]byte{4, 5}); err != nil {
		fmt.Printf("2 write fail: %v\n", err)
		return
	}
	if want, have := uint64(5), atomic.LoadUint64(&dialConn.wr); want != have {
		fmt.Printf("want: %d, have: %d\n", want, have)
		return
	}

	dialConn, dialErr = nil, errors.New("nono")
	mgr.Put(errors.New("trigger that reconnect y'all"))
	if conn = mgr.Take(); conn != nil {
		fmt.Printf("want nil conn, got real conn\n")
		return
	}

	go func() {
		done := time.After(100 * time.Millisecond)
		for {
			select {
			case tickChan <- time.Now():
			case <-done:
				return
			}
		}
	}()

	if within(100*time.Millisecond, func() bool {
		conn = mgr.Take()
		return conn != nil
	}) {
		fmt.Printf("eventually got a good conn, despite failing dialer\n")
		return
	}
}

type mockConn struct {
	rd, wr uint64
}

func (c *mockConn) Read(b []byte) (n int, err error) {
	atomic.AddUint64(&c.rd, uint64(len(b)))
	return len(b), nil
}

func (c *mockConn) Write(b []byte) (n int, err error) {
	atomic.AddUint64(&c.wr, uint64(len(b)))
	return len(b), nil
}

func (c *mockConn) Close() error                       { return nil }
func (c *mockConn) LocalAddr() net.Addr                { return nil }
func (c *mockConn) RemoteAddr() net.Addr               { return nil }
func (c *mockConn) SetDeadline(t time.Time) error      { return nil }
func (c *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (c *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func within(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for {
		if time.Now().After(deadline) {
			return false
		}
		if f() {
			return true
		}
		time.Sleep(d / 10)
	}
}
