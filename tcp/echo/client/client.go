package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/pyihe/go-pkg/packets"
	"github.com/pyihe/go-pkg/strings"
	"github.com/pyihe/go-pkg/syncs"
)

var (
	pkt = packets.NewPacket(4, 4096)
)

func main() {
	wg := &syncs.WgWrapper{}
	for i := 0; i < 1; i++ {
		wg.Wrap(func() {
			dialTCP("localhost:9999").start()
			//dialWithTLS("localhost:9999", "../certs/client.pem", "../certs/client.key").start()
		})
	}
	wg.Wait()
}

type client struct {
	conn net.Conn
	w    *bufio.Writer
}

func dialTCP(addr string) (c *client) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	c = &client{conn: conn, w: bufio.NewWriter(conn)}
	c.onConnect()
	return
}

func dialWithTLS(addr string, cert, key string) (c *client) {
	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		panic(err)
	}

	cCertBytes, err := ioutil.ReadFile(cert)
	if err != nil {
		panic(err)
	}
	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM(cCertBytes); !ok {
		panic("failed to append client pem")
	}

	conf := &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		ClientCAs:          clientCertPool,
		InsecureSkipVerify: true, // 这里设置为true表示不需要验证服务器主机名与证书主机名是否一致，只当测试时设置
	}
	conn, err := tls.Dial("tcp", addr, conf)
	if err != nil {
		panic(err)
	}
	c = &client{
		conn: conn,
		w:    bufio.NewWriter(conn),
	}
	c.onConnect()
	return
}

func (c *client) sendMsg(message []byte) error {
	data, err := pkt.Packet(message)
	if err != nil {
		return err
	}
	_, err = c.w.Write(data)
	if err != nil {
		return err
	}
	return c.w.Flush()
}

func (c *client) start() {
	reader := bufio.NewReader(c.conn)
	for {
		data, err := pkt.UnPacket(reader)
		if err != nil {
			fmt.Println("dddd", err)
			break
		}
		c.onMessage(data)
	}
}

// 建立连接时执行
func (c *client) onConnect() {
	c.sendMsg(strings.Bytes("Welcome to Go World!"))
}

// 收到消息时执行: 这里只是简单的将消息
func (c *client) onMessage(message []byte) {
	fmt.Println(string(message))
	c.sendMsg(message)
}
