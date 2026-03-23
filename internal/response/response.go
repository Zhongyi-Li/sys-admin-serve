package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess         = 0
	CodeBadRequest      = 4000
	CodeUnauthorized    = 4001
	CodeForbidden       = 4003
	CodeNotFound        = 4004
	CodeInternalError   = 5000
	defaultSuccessMsg   = "success"
	defaultErrorMessage = "internal server error"
)

type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type PageData struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

func Success(c *gin.Context, data any) {
	JSON(c, http.StatusOK, CodeSuccess, defaultSuccessMsg, data)
}

func SuccessWithMessage(c *gin.Context, message string, data any) {
	JSON(c, http.StatusOK, CodeSuccess, message, data)
}

func Page(c *gin.Context, list any, total int64, page, pageSize int) {
	Success(c, PageData{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func BadRequest(c *gin.Context, message string) {
	JSON(c, http.StatusBadRequest, CodeBadRequest, message, nil)
}

func Unauthorized(c *gin.Context, message string) {
	JSON(c, http.StatusUnauthorized, CodeUnauthorized, message, nil)
}

func Forbidden(c *gin.Context, message string) {
	JSON(c, http.StatusForbidden, CodeForbidden, message, nil)
}

func NotFound(c *gin.Context, message string) {
	JSON(c, http.StatusNotFound, CodeNotFound, message, nil)
}

func InternalError(c *gin.Context) {
	JSON(c, http.StatusInternalServerError, CodeInternalError, defaultErrorMessage, nil)
}

func JSON(c *gin.Context, httpStatus, code int, message string, data any) {
	c.JSON(httpStatus, Body{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
