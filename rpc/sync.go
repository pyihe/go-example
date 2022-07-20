package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"time"

	"github.com/pyihe/go-pkg/errors"
	"github.com/pyihe/go-pkg/syncs"
)

//func main() {
//	s, err := serveHTTP(":8888")
//	if err != nil {
//		panic(err)
//	}
//	if err = s.Register(new(Game)); err != nil {
//		panic(err)
//	}
//	runClient(true)
//}

func serveTCP(addr string) (*rpc.Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	go func() {
		rpc.Accept(ln)
	}()
	return rpc.DefaultServer, nil
}

func serveHTTP(addr string) (*rpc.Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	s := rpc.NewServer()
	rpc.HandleHTTP()
	go func() {
		http.Serve(ln, s)
	}()
	return s, nil
}

func runClient(isHttp bool) {
	var client *rpc.Client
	var err error
	if isHttp {
		client, err = rpc.DialHTTP("tcp", ":8888")
		if err != nil {
			fmt.Println("dial err: ", err)
			return
		}
	} else {
		client, err = rpc.Dial("tcp", ":8888")
		if err != nil {
			fmt.Println("dial err: ", err)
			return
		}
	}

	v := new(int64)
	if err = client.Call("Game.Login", 1, v); err != nil {
		fmt.Println("call err: ", err)
		return
	}
	fmt.Println(*v)
	time.Sleep(2 * time.Second)
	if err = client.Call("Game.Logoff", 1, v); err != nil {
		fmt.Println("call err1: ", err)
	}
	time.Sleep(2 * time.Second)
	fmt.Println(*v)
	if err = client.Call("Game.Login", 1, v); err != nil {
		fmt.Println("call err2: ", err)
		return
	}
	fmt.Println(*v)
}

type Game struct {
	history sync.Map
	players sync.Map
	counter syncs.AtomicInt64
}

// Login 玩家登录，传入玩家ID，返回上次登录时间
func (g *Game) Login(playerId int64, history *int64) error {
	if _, online := g.players.Load(playerId); online {
		return errors.New("重复登录!")
	}
	// 记录登录信息
	g.players.Store(playerId, time.Now().Unix())
	// 增加在线人数
	g.counter.Inc(1)
	// 返回上次登录时间
	if v, ok := g.history.Load(playerId); ok {
		*history = v.(int64)
	}
	return nil
}

// Logoff 登出, 返回在线时长，单位秒
func (g *Game) Logoff(playerId int64, sec *int64) error {
	loginAt, online := g.players.Load(playerId)
	if !online {
		return errors.New("Not Found")
	}
	nowUnix := time.Now().Unix()
	// 保存本次登录历史
	g.history.Store(playerId, nowUnix)
	// 删除登录信息
	g.players.Delete(playerId)
	// 减少在线数量
	g.counter.Inc(-1)
	*sec = nowUnix - loginAt.(int64)
	return nil
}
