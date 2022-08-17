package route

import (
	"github.com/gin-gonic/gin"
	"github.com/pyihe/go-example/gin-swagger/model/form"
	"github.com/pyihe/go-example/gin-swagger/service"
	"github.com/pyihe/go-pkg/https/http_api"
)

type departmentRouter struct {
	service *service.DepartmentService
}

func NewDepartmentRouter(s *service.DepartmentService) http_api.APIHandler {
	return &departmentRouter{service: s}
}

func (d *departmentRouter) Handle(r http_api.IRouter) {
	group := r.Group("/depart")
	{
		group.POST("", http_api.WrapHandler(d.Add))
		group.DELETE("/:name", http_api.WrapHandler(d.Delete))
		group.PATCH("/:name", http_api.WrapHandler(d.Modify))
		group.GET("/:name", http_api.WrapHandler(d.Information))
	}
}

// Add 添加部门
// @Summary 		添加部门
// @Description 	添加部门信息
// @Tags 			Department
// @Accept 			json
// @Produce 		json
// @Param 			addBody 	body 		form.AddRequest true "部门名称等信息"
// @Success 		200
// @Failure 		400
// @Router 			/api/depart 		[POST]
func (d *departmentRouter) Add(c *gin.Context) (result interface{}, err error) {
	var req form.AddRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}
	err = d.service.Add(req.Name)
	return
}

// Delete 删除部门信息
// @Summary 		删除部门
// @Description 	删除部门信息
// @Tags 			Department
// @Accept 			json
// @Produce 		json
// @Param 			name 	path 		string true "部门名称"
// @Success 		200
// @Failure 		400
// @Router 			/api/depart/{name} 		[DELETE]
func (d *departmentRouter) Delete(c *gin.Context) (result interface{}, err error) {
	err = d.service.Delete(c.Param("name"))
	return
}

// Modify 修改部门信息
// @Summary 		修改部门
// @Description 	修改部门信息
// @Tags 			Department
// @Accept 			json
// @Produce 		json
// @Param 			name 	path 		string true "部门名称"
// @Success 		200
// @Failure 		400
// @Router 			/api/depart/{name} 		[PATCH]
func (d *departmentRouter) Modify(c *gin.Context) (result interface{}, err error) {
	err = d.service.Modify(c.Param("name"))
	return
}

// Information 查询部门信息
// @Summary 		查询部门
// @Description 	查询部门信息
// @Tags 			Department
// @Accept 			json
// @Produce 		json
// @Param 			name 	path 		string true "部门名称"
// @Success 		200
// @Failure 		400
// @Router 			/api/depart/{name} 		[GET]
func (d *departmentRouter) Information(c *gin.Context) (result interface{}, err error) {
	return d.service.Information(c.Param("name"))
}
