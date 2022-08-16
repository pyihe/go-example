package route

import (
	"github.com/gin-gonic/gin"
	"github.com/pyihe/go-example/gin-swagger/form"
	"github.com/pyihe/go-example/gin-swagger/service"
	"github.com/pyihe/go-pkg/https/http_api"
	"github.com/pyihe/plogs"
)

type userRoute struct {
	service *service.UserService
}

func NewUserRoute(s *service.UserService) http_api.APIHandler {
	return &userRoute{
		service: s,
	}
}

func (u *userRoute) Handle(r http_api.IRouter) {
	r.POST("/login", http_api.WrapFunc(u.Login))
	r.POST("/login1", u.Login1)
}

// Login
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Router       /login1?user_name=?&password=? [post]
func (u *userRoute) Login(c *gin.Context) (result interface{}, err error) {
	var req form.LoginRequest
	if err = c.ShouldBind(&req); err != nil {
		return
	}

	return u.service.Login(req.UserName, req.Password)
}

// Login1
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Router       /login1?user_name=?&password=? [post]
func (u *userRoute) Login1(c *gin.Context) {
	var req form.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		plogs.Errorf("bind发生错误: %v", err)
		return
	}

	result, err := u.service.Login(req.UserName, req.Password)
	if err != nil {
		plogs.Errorf("登录出错: %v", err)
		return
	}
	http_api.IndentedJSON(c, err, result)
	return
}
