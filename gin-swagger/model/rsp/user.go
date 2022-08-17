package rsp

// Login 登录回复
type Login struct {
	Username string `json:"username,omitempty"` // 用户名
	Token    string `json:"token,omitempty"`    // Token令牌
}

// Register 注册回复
type Register struct {
	Id   int64  `form:"id"          json:"id"`   // 用户ID
	Name string `form:"name"        json:"name"` // 姓名
	Sex  uint8  `form:"sex"         json:"sex"`  // 性别
	Age  int    `form:"age"         json:"age"`  // 年龄
	City string `form:"city"        json:"city"` // 城市
}

// Information 用户信息
type Information struct {
	Id   int64  `form:"id"          json:"id"`   // 用户ID
	Name string `form:"name"        json:"name"` // 姓名
	Sex  uint8  `form:"sex"         json:"sex"`  // 性别
	Age  int    `form:"age"         json:"age"`  // 年龄
	City string `form:"city"        json:"city"` // 城市
}
