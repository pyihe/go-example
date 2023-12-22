package httpnet

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pyihe/go-pkg/errors"
)

const (
	success       = "SUCCESS"
	Authorization = "Authorization"
	TokenKey      = "client_id"
)

type response struct {
	Code    int32       `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func IndentedJSON(c *gin.Context, err error, data interface{}) {
	var rsp = &response{}
	var status = http.StatusOK

	switch {
	case err != nil:
		status = http.StatusBadRequest
		switch err.(type) {
		case *errors.Error:
			e := err.(*errors.Error)
			rsp.Code = e.Code()
			rsp.Message = e.Message()
		default:
			rsp.Message = err.Error()
		}
	default:
		rsp.Message = success
		rsp.Data = data
	}
	c.IndentedJSON(status, rsp)
}

func WrapHandler(handler func(c *gin.Context) (interface{}, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		if handler == nil {
			return
		}
		data, err := handler(c)
		IndentedJSON(c, err, data)
	}
}
