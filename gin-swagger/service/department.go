package service

import (
	"sync"

	"github.com/pyihe/go-pkg/errors"
)

var (
	errDepartmentExist = errors.New("部门已存在!")
)

type DepartmentService struct {
	mu          sync.RWMutex
	departments map[string]struct{}
}

func NewDepartmentService() *DepartmentService {
	return &DepartmentService{
		mu:          sync.RWMutex{},
		departments: make(map[string]struct{}),
	}
}

// Add 添加部门
func (d *DepartmentService) Add(name string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	_, exist := d.departments[name]
	if exist {
		return errDepartmentExist
	}
	d.departments[name] = struct{}{}
	return nil
}

// Delete 删除部门信息
func (d *DepartmentService) Delete(name string) error {
	d.mu.Lock()
	delete(d.departments, name)
	d.mu.Unlock()
	return nil
}

// Modify 修改部门信息
func (d *DepartmentService) Modify(name string) error {
	return nil
}

// Information 查询部门信息
func (d *DepartmentService) Information(name string) (interface{}, error) {
	type depart struct {
		Name string `form:"name"         json:"name"`
	}
	return &depart{Name: name}, nil
}
