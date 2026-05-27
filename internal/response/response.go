package response

import "github.com/gin-gonic/gin"

const (
	CodeOK             = 0
	CodeInvalidRequest = 40001
	CodeUnauthorized   = 40101
	CodeForbidden      = 40301
	CodeNotFound       = 40401
	CodeInternalError  = 50001
)

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(200, Body{Code: CodeOK, Message: "ok", Data: data})
}

func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Body{Code: code, Message: message, Data: nil})
}
