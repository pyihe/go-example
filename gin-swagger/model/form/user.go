package form

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `form:"username"  json:"username"` // 用户名
	Password string `form:"password"  json:"password"` // 密码
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `form:"username"         json:"username"` // 用户名
	Password string `form:"password"         json:"password"` // 密码
	Sex      uint8  `form:"sex"              json:"sex"`      // 性别
	Age      int    `form:"age"              json:"age"`      // 年龄
	City     string `form:"city"             json:"city"`     // 城市
}

// ModifyRequest 修改用户信息请求
type ModifyRequest struct {
	Sex  uint8  `form:"sex"         json:"sex"`  // 性别
	Age  int    `form:"age"         json:"age"`  // 年龄
	City string `form:"city"        json:"city"` // 城市
}
