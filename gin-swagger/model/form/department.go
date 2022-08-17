package form

// AddRequest 添加部门的请求
type AddRequest struct {
	Name string `form:"name"         json:"name"` // 部门名称
}
