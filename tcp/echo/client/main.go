package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"

	"github.com/pyihe/go-pkg/packets"
	"github.com/pyihe/go-pkg/strings"
	"github.com/pyihe/go-pkg/syncs"
)

var (
	opts = []packets.Option{
		packets.WithHeaderSize(4),
		packets.WithMaxMsgSize(4 * 1024),
		packets.WithMinMsgSize(1),
	}
	pkt = packets.NewPacket(opts...)
)

func main() {
	wg := &syncs.WgWrapper{}
	for i := 0; i < 1; i++ {
		wg.Wrap(func() {
			dialTCP("127.0.0.1:9999").start()
			//dialWithTLS("localhost:9999", "cadir", "../certs/client.pem", "../certs/client.key").start()
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

func dialWithTLS(addr string, ca, cert, key string) (c *client) {
	// 加载客户端证书
	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		panic(err)
	}
	// 因为是自签名证书, 这里直接使用的是客户端证书, 实际上如果需要执行完整的证书链验证过程: 验证服务端主机名与证书中的是否一致
	//cCertBytes, err := ioutil.ReadFile(ca)
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
		ClientCAs:          clientCertPool, // 用于验证服务端证书是否合法
		InsecureSkipVerify: true,           // 这里设置为true表示不需要验证服务器主机名与证书主机名是否一致，只当测试时设置
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
			break
		}
		c.onMessage(data)
	}
}

// 建立连接时执行
func (c *client) onConnect() {
	c.sendMsg(strings.Bytes("Welcome to the Gopher's World!"))
}

// 收到消息时执行: 这里只是简单的将消息
func (c *client) onMessage(message []byte) {
	c.sendMsg(message)
}
