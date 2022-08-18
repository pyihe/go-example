package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pyihe/go-example/gin-swagger/model/form"
	"github.com/pyihe/go-example/gin-swagger/service"
	"github.com/pyihe/go-example/pkg"
	"github.com/pyihe/go-pkg/https/http_api"
)

type userRouter struct {
	service *service.UserService
}

func NewUserRoute(s *service.UserService) http_api.APIHandler {
	return &userRouter{
		service: s,
	}
}

func (u *userRouter) Handle(r http_api.IRouter) {
	group := r.Group("/user")

	unCheck := group.Group("")
	{
		unCheck.POST("/login", http_api.WrapHandler(u.Login))
		unCheck.POST("/register", http_api.WrapHandler(u.Register))
	}

	auth := r.Group("/user", u.auth())
	{
		auth.DELETE("/:username", http_api.WrapHandler(u.Delete))
		auth.PATCH("/:username", http_api.WrapHandler(u.Modify))
		auth.GET("/:username", http_api.WrapHandler(u.Information))
		auth.GET("", http_api.WrapHandler(u.List))
	}
}

// Login 			登录
// @Summary 		登录API
// @Description 	用户通过用户名/密码登录系统
// @Tags 			User
// @Accept 			json
// @Produce 		json
// @Param 			loginBody 	body 		form.LoginRequest true "用户名"
// @Success 		200 				{object}	rsp.Login
// @Failure 		400
// @Router 			/api/user/login 				[POST]
func (u *userRouter) Login(c *gin.Context) (result interface{}, err error) {
	var req form.LoginRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}

	return u.service.Login(req.Username, req.Password)
}

// Register 		注册
// @Summary 		注册API
// @Description 	用户注册
// @Tags 			User
// @Accept 			json
// @Produce 		json
// @Param 			registerBody 	body 		form.RegisterRequest true "用户名"
// @Success 		200 	{object}	rsp.Register
// @Failure 		400
// @Router 			/api/user/register 	[POST]
func (u *userRouter) Register(c *gin.Context) (result interface{}, err error) {
	var req form.RegisterRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}

	return u.service.Register(&req)
}

// Delete 删除用户
// @Summary 		删除API
// @Description 	删除用户
// @Tags 			User
// @Produce 		json
// @Param 			username path string true "用户名"
// @Security 		ApiKeyAuth
// @Success 		200
// @Failure 		400
// @Router 			/api/user/{username} 	[DELETE]
func (u *userRouter) Delete(c *gin.Context) (result interface{}, err error) {
	err = u.service.Delete(c.Param("username"))
	return
}

// Modify 修改用户信息
// @Summary 		修改API
// @Description 	修改用户信息
// @Tags 			User
// @Accept 			json
// @Produce 		json
// @Param 			modifyBody body form.ModifyRequest true "用户名"
// @Security 		ApiKeyAuth
// @Success 		200
// @Failure 		400
// @Router 			/api/user/{username} 	[PATCH]
func (u *userRouter) Modify(c *gin.Context) (result interface{}, err error) {
	var req form.ModifyRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}
	err = u.service.Modify(c.Param("username"), &req)
	return
}

// Information 用户信息
// @Summary 		查询API
// @Description 	查询用户信息
// @Tags 			User
// @Produce 		json
// @Param 			username path string true "用户名"
// @Security 		ApiKeyAuth
// @Success 		200			{object} 	rsp.Information
// @Failure 		400
// @Router 			/api/user/{username} 	[GET]
func (u *userRouter) Information(c *gin.Context) (result interface{}, err error) {
	return u.service.Information(c.Param("username"))
}

// List 用户列表
// @Summary 		获取用户列表
// @Description 	获取用户列表
// @Tags 			User
// @Produce 		json
// @Security 		ApiKeyAuth
// @Success 		200
// @Failure 		400
// @Router 			/api/user 	[GET]
func (u *userRouter) List(c *gin.Context) (result interface{}, err error) {
	return u.service.List()
}

func (u *userRouter) auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(pkg.Authorization)
		user, err := u.service.GetUserByToken(token)
		if err != nil {
			http_api.IndentedJSON(c, err, nil)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("userId", user.Id)
		c.Next()
	}
}
