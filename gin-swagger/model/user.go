package model

// User 用户信息
type User struct {
	Id       int64  `form:"id"          json:"id"`   // 用户ID
	Name     string `form:"name"        json:"name"` // 姓名
	Sex      uint8  `form:"sex"         json:"sex"`  // 性别
	Age      int    `form:"age"         json:"age"`  // 年龄
	City     string `form:"city"        json:"city"` // 城市
	Password string `form:"-"           json:"-"`    // 密码
	Token    string `form:"-"           json:"-"`    // token
}
