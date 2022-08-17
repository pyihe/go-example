package service

import (
	"sync"

	"github.com/pyihe/go-example/gin-swagger/model"
	"github.com/pyihe/go-example/gin-swagger/model/form"
	"github.com/pyihe/go-example/gin-swagger/model/rsp"
	"github.com/pyihe/go-pkg/errors"
	"github.com/pyihe/go-pkg/rands"
	"github.com/pyihe/go-pkg/snowflakes"
	"github.com/pyihe/plogs"
)

var (
	errWrongPass             = errors.New("密码错误!")
	errUserNotExist          = errors.New("用户不存在!")
	errInvalidUsernameFormat = errors.New("用户名格式错误!")
	errUsernameExist         = errors.New("用户名已存在!")
	errInvalidAuthorization  = errors.New("非法的验证信息!")
	errInvalidPassFormat     = errors.New("密码格式错误!")
)

type UserService struct {
	mu        sync.RWMutex
	idFactory snowflakes.Worker      // 生成ID
	accounts  map[string]*model.User // <用户名, 用户信息>
	tokens    map[string]string      // <token, 用户名>
}

func NewUserService() *UserService {
	s := &UserService{}
	s.idFactory = snowflakes.NewWorker(1)
	s.accounts = map[string]*model.User{
		"admin": &model.User{
			Id:       1,
			Name:     "admin",
			Sex:      1,
			Age:      20,
			City:     "ChengDu",
			Password: "adminadmin",
			Token:    "",
		},
	}
	s.tokens = make(map[string]string)
	return s
}

// GetUserByToken 根据token获取用户信息
func (u *UserService) GetUserByToken(token string) (*model.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	username, exist := u.tokens[token]
	if !exist {
		return nil, errInvalidAuthorization
	}
	user := u.accounts[username]
	return user, nil
}

// Login 账号密码登录，返回Token凭证
func (u *UserService) Login(username, pass string) (response rsp.Login, err error) {
	plogs.Infof("玩家[%v]请求登录, 密码为: [%v]", username, pass)
	u.mu.RLock()
	defer u.mu.RUnlock()

	user, exist := u.accounts[username]
	if !exist {
		err = errUserNotExist
		return
	}

	if user.Password != pass {
		err = errWrongPass
		return
	}

	response.Username = username
	response.Token = rands.String(32)
	u.tokens[response.Token] = username
	user.Token = response.Token
	plogs.Infof("玩家[%v]登录成功, 返回Token: [%s]", username, response.Token)
	return
}

// Register 用户注册
func (u *UserService) Register(info *form.RegisterRequest) (*rsp.Register, error) {
	plogs.Infof("收到注册请求: %+v", *info)
	// 校验用户名
	if len(info.Username) < 4 {
		return nil, errInvalidUsernameFormat
	}

	// 校验密码
	if len(info.Password) < 8 {
		return nil, errInvalidPassFormat
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	_, exist := u.accounts[info.Username]
	if exist {
		return nil, errUsernameExist
	}

	user := &model.User{
		Id:       u.idFactory.GetInt64(),
		Name:     info.Username,
		Sex:      info.Sex,
		Age:      info.Age,
		City:     info.City,
		Password: info.Password,
	}
	u.accounts[user.Name] = user

	response := &rsp.Register{
		Id:   user.Id,
		Name: user.Name,
		Sex:  user.Sex,
		Age:  user.Age,
		City: user.City,
	}
	plogs.Infof("玩家[%v]注册成功!", user.Name)
	return response, nil
}

// Delete 删除用户
func (u *UserService) Delete(username string) error {
	plogs.Infof("收到删除玩家的请求: [%v]", username)
	if len(username) == 0 {
		return nil
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	user, exist := u.accounts[username]
	if !exist {
		return nil
	}
	delete(u.accounts, username)
	delete(u.tokens, user.Token)
	plogs.Infof("玩家[%v]删除成功", username)
	return nil
}

// Modify 修改用户信息
func (u *UserService) Modify(username string, info *form.ModifyRequest) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	user, exist := u.accounts[username]
	if !exist {
		return errUserNotExist
	}

	user.Age = info.Age
	user.Sex = info.Sex
	user.City = info.City
	plogs.Infof("删除玩家[%v]信息成功!", username)
	return nil
}

// Information 获取用户个人信息
func (u *UserService) Information(username string) (*rsp.Information, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	user, exist := u.accounts[username]
	if !exist {
		return nil, errUserNotExist
	}

	info := &rsp.Information{
		Id:   user.Id,
		Name: user.Name,
		Sex:  user.Sex,
		Age:  user.Age,
		City: user.City,
	}
	plogs.Infof("获取玩家[%v]信息成功", username)
	return info, nil
}

// List 获取用户列表
func (u *UserService) List() ([]*rsp.Information, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	var list = make([]*rsp.Information, 0, len(u.accounts))
	for _, user := range u.accounts {
		info := &rsp.Information{
			Id:   user.Id,
			Name: user.Name,
			Sex:  user.Sex,
			Age:  user.Age,
			City: user.City,
		}
		list = append(list, info)
	}
	plogs.Infof("获取玩家列表成功!")
	return list, nil
}
