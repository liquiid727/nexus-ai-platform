package errors

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

// usage
// err := new (ErrorX)
// err.WithMessage("Message").WithMedata(map[string]string{"age": "0"})

type ErrorX struct {
	// 表示错误的HTTP状态码，例如400, 500等
	Code int `json:"code,omitempty"`
	// 错误原因,业务错误，用于精准定位问题
	Reason string `json:"reason"`
	// 错误信息,可以暴露给用户看
	Message string `json:"message"`

	// Metadata 额外的元数据，用于记录错误的上下文信息
	Metadata map[string]string `json:"metadata,omitempty"`
}

// 支持New函数
func New(code int, reason string, format string, args ...any) *ErrorX {
	return &ErrorX{
		Code:    code,
		Reason:  reason,
		Message: fmt.Sprintf(format, args...),
	}
}

// Error 实现error接口，返回错误信息
func (Err *ErrorX) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v", Err.Code, Err.Reason, Err.Message, Err.Metadata)
}

// WithMessage 设置错误信息
func (Err *ErrorX) WithMessage(format string, args ...any) *ErrorX {
	Err.Message = fmt.Sprintf(format, args...)
	return Err
}

// WithMetadata 设置额外的元数据
func (Err *ErrorX) WithMetadata(md map[string]string) *ErrorX {
	Err.Metadata = md
	return Err
}

// 具体的使用key-value对设置元数据
func (Err *ErrorX) KV(kvs ...string) *ErrorX {
	if Err.Metadata == nil {
		Err.Metadata = make(map[string]string)
	}
	for i := 0; i < len(kvs); i += 2 {
		if i+1 >= len(kvs) {
			break
		}
		Err.Metadata[kvs[i]] = kvs[i+1]
	}
	return Err
}

// 满足gRPC错误接口
func (Err *ErrorX) GRPCStatus() *status.Status {
	// TODO
	// return status.New(Err.Code, Err.Message)
	return nil
}

// WithRequestID 设置请求ID
func (Err *ErrorX) WithRequestID(requestID string) *ErrorX {
	return Err.KV("X-nexus-gateway-request_id", requestID)
}

// Is 实现error接口，判断是否是相同的错误
// 递归错误链，比较ErrorX实例的Code和Reason字段
func (Err *ErrorX) Is(target error) bool {
	if errx := new(ErrorX); errors.As(target, &errx) {
		return errx.Code == Err.Code && errx.Reason == Err.Reason
	}
	return false
}

func Code(err error) int {
	if err == nil {
		return http.StatusOK
	}
	return FromError(err).Code
}
func Reason(err error) string {
	if err == nil {
		// 	return ErrInternal.Reason
		return " "
	}
	return FromError(err).Reason
}

// 将err转化为 自定义的*ErrorX
func FromError(err error) *ErrorX {
	// 传入的为nil，表示没有错误需要处理
	if err == nil {
		return nil
	}
	// 核对error，是否以及是ErrorX类型的错误
	// 如果是的，则直接返回
	if errx := new(ErrorX); errors.As(err, &errx) {
		return errx
	}

	// 根据错误类型进行返回，是GPRC还是HTTP
	if _, ok := status.FromError(err); !ok {
		return New(ErrInternal.Code, ErrInternal.Reason, "%s", err.Error())
	}
	// TODO
	// 如果是HTTP错误，则返回HTTP错误
	if s, ok := status.FromError(err); ok {
		ret := New(int(s.Code()), s.Message(), "%s", err.Error())
		for _, detail := range s.Details() {
			if typed, ok := detail.(*errdetails.ErrorInfo); ok {
				ret.Reason = typed.Reason
				return ret.WithMetadata(typed.Metadata)
			}
		}

	}
	return nil
}
