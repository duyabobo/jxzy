package errorx

import "fmt"

// 错误代码常量
const (
	// 通用错误
	ErrCodeSuccess      = 200
	ErrCodeParamError   = 400
	ErrCodeUnauthorized = 401
	ErrCodeForbidden    = 403
	ErrCodeNotFound     = 404
	ErrCodeSystemError  = 500

	// 业务错误代码 (10000-19999)
	ErrCodeUserNotFound    = 10001
	ErrCodeUserExists      = 10002
	ErrCodeInvalidPassword = 10003
	ErrCodeTokenExpired    = 10004
	ErrCodeTokenInvalid    = 10005

	// 上下文错误 (11000-11999)
	ErrCodeContextNotFound    = 11001
	ErrCodeContextExists      = 11002
	ErrCodeContextNotBelongTo = 11003

	// Prompt错误 (12000-12999)
	ErrCodePromptNotFound     = 12001
	ErrCodePromptExists       = 12002
	ErrCodePromptInvalid      = 12003
	ErrCodePromptNotBelongTo  = 12004
	ErrCodePromptRenderFailed = 12005

	// LLM错误 (13000-13999)
	ErrCodeLLMNotAvailable  = 13001
	ErrCodeLLMQuotaExceeded = 13002
	ErrCodeLLMRequestFailed = 13003
	ErrCodeLLMTimeout       = 13004

	// RAG错误 (14000-14999)
	ErrCodeDocumentNotFound      = 14001
	ErrCodeDocumentUploadFailed  = 14002
	ErrCodeDocumentProcessFailed = 14003
	ErrCodeCollectionNotFound    = 14004
	ErrCodeVectorSearchFailed    = 14005
)

// CodeError 业务错误结构
type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("Code: %d, Msg: %s", e.Code, e.Msg)
}

func (e *CodeError) GetCode() int {
	return e.Code
}

func (e *CodeError) GetMsg() string {
	return e.Msg
}

// 预定义错误
var (
	ErrSuccess      = &CodeError{Code: ErrCodeSuccess, Msg: "success"}
	ErrParamError   = &CodeError{Code: ErrCodeParamError, Msg: "参数错误"}
	ErrUnauthorized = &CodeError{Code: ErrCodeUnauthorized, Msg: "未授权"}
	ErrForbidden    = &CodeError{Code: ErrCodeForbidden, Msg: "禁止访问"}
	ErrNotFound     = &CodeError{Code: ErrCodeNotFound, Msg: "资源不存在"}
	ErrSystemError  = &CodeError{Code: ErrCodeSystemError, Msg: "系统错误"}

	// 用户相关错误
	ErrUserNotFound    = &CodeError{Code: ErrCodeUserNotFound, Msg: "用户不存在"}
	ErrUserExists      = &CodeError{Code: ErrCodeUserExists, Msg: "用户已存在"}
	ErrInvalidPassword = &CodeError{Code: ErrCodeInvalidPassword, Msg: "密码错误"}
	ErrTokenExpired    = &CodeError{Code: ErrCodeTokenExpired, Msg: "Token已过期"}
	ErrTokenInvalid    = &CodeError{Code: ErrCodeTokenInvalid, Msg: "Token无效"}

	// 上下文相关错误
	ErrContextNotFound    = &CodeError{Code: ErrCodeContextNotFound, Msg: "上下文不存在"}
	ErrContextExists      = &CodeError{Code: ErrCodeContextExists, Msg: "上下文已存在"}
	ErrContextNotBelongTo = &CodeError{Code: ErrCodeContextNotBelongTo, Msg: "上下文不属于当前用户"}

	// Prompt相关错误
	ErrPromptNotFound     = &CodeError{Code: ErrCodePromptNotFound, Msg: "Prompt不存在"}
	ErrPromptExists       = &CodeError{Code: ErrCodePromptExists, Msg: "Prompt已存在"}
	ErrPromptInvalid      = &CodeError{Code: ErrCodePromptInvalid, Msg: "Prompt格式无效"}
	ErrPromptNotBelongTo  = &CodeError{Code: ErrCodePromptNotBelongTo, Msg: "Prompt不属于当前用户"}
	ErrPromptRenderFailed = &CodeError{Code: ErrCodePromptRenderFailed, Msg: "Prompt渲染失败"}

	// LLM相关错误
	ErrLLMNotAvailable  = &CodeError{Code: ErrCodeLLMNotAvailable, Msg: "LLM服务不可用"}
	ErrLLMQuotaExceeded = &CodeError{Code: ErrCodeLLMQuotaExceeded, Msg: "LLM配额已用尽"}
	ErrLLMRequestFailed = &CodeError{Code: ErrCodeLLMRequestFailed, Msg: "LLM请求失败"}
	ErrLLMTimeout       = &CodeError{Code: ErrCodeLLMTimeout, Msg: "LLM请求超时"}

	// RAG相关错误
	ErrDocumentNotFound      = &CodeError{Code: ErrCodeDocumentNotFound, Msg: "文档不存在"}
	ErrDocumentUploadFailed  = &CodeError{Code: ErrCodeDocumentUploadFailed, Msg: "文档上传失败"}
	ErrDocumentProcessFailed = &CodeError{Code: ErrCodeDocumentProcessFailed, Msg: "文档处理失败"}
	ErrCollectionNotFound    = &CodeError{Code: ErrCodeCollectionNotFound, Msg: "文档集合不存在"}
	ErrVectorSearchFailed    = &CodeError{Code: ErrCodeVectorSearchFailed, Msg: "向量搜索失败"}
)

// NewCodeError 创建新的业务错误
func NewCodeError(code int, msg string) *CodeError {
	return &CodeError{Code: code, Msg: msg}
}

// NewCodeErrorf 创建新的业务错误（支持格式化）
func NewCodeErrorf(code int, format string, args ...interface{}) *CodeError {
	return &CodeError{Code: code, Msg: fmt.Sprintf(format, args...)}
}
