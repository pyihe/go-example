package form

type LoginRequest struct {
	UserName string `form:"user_name" json:"user_name"`
	Password string `form:"password"  json:"password"`
}
