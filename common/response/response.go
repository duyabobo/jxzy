package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 成功响应
func Success(w http.ResponseWriter, data interface{}) {
	httpx.OkJson(w, &Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// 错误响应
func Error(w http.ResponseWriter, code int, message string) {
	httpx.WriteJson(w, code, &Response{
		Code:    code,
		Message: message,
	})
}

// 业务错误响应
func BusinessError(w http.ResponseWriter, message string) {
	httpx.WriteJson(w, http.StatusOK, &Response{
		Code:    400,
		Message: message,
	})
}

// 系统错误响应
func SystemError(w http.ResponseWriter, message string) {
	httpx.WriteJson(w, http.StatusInternalServerError, &Response{
		Code:    500,
		Message: message,
	})
}

// 参数错误响应
func ParamError(w http.ResponseWriter, message string) {
	httpx.WriteJson(w, http.StatusBadRequest, &Response{
		Code:    400,
		Message: message,
	})
}

// 未授权响应
func Unauthorized(w http.ResponseWriter, message string) {
	httpx.WriteJson(w, http.StatusUnauthorized, &Response{
		Code:    401,
		Message: message,
	})
}

// 禁止访问响应
func Forbidden(w http.ResponseWriter, message string) {
	httpx.WriteJson(w, http.StatusForbidden, &Response{
		Code:    403,
		Message: message,
	})
}

// 资源不存在响应
func NotFound(w http.ResponseWriter, message string) {
	httpx.WriteJson(w, http.StatusNotFound, &Response{
		Code:    404,
		Message: message,
	})
}
