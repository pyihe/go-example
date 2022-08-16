package service

import (
	"sync"

	gin_swagger "github.com/pyihe/go-example/gin-swagger"
	"github.com/pyihe/go-pkg/errors"
	"github.com/pyihe/go-pkg/rands"
	"github.com/pyihe/plogs"
)

type UserService struct {
	mu       sync.RWMutex
	accounts map[string]string // <用户名, 密码>
	tokens   map[string]string // <token, 用户名>
}

func NewUserService() *UserService {
	s := &UserService{}
	s.accounts = map[string]string{
		"admin": "admin",
		"user1": "user1",
	}
	s.tokens = make(map[string]string)
	return s
}

// Login 账号密码登录，返回Token凭证
func (u *UserService) Login(user, pass string) (response gin_swagger.LoginResponse, err error) {
	plogs.Infof("玩家[%v]请求登录, 密码为: [%v]", user, pass)
	u.mu.RLock()
	defer u.mu.RUnlock()

	password, exist := u.accounts[user]
	if !exist {
		err = errors.New("用户不存在")
		return
	}

	if password != pass {
		err = errors.New("密码错误")
		return
	}

	response.Username = user
	response.Token = rands.String(32)
	u.tokens[response.Token] = user
	plogs.Infof("玩家[%v]登录成功, 返回Token: [%s]", user, response.Token)
	return
}
