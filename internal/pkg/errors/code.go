package errors

import "net/http"

var (
	// 代表请求成功
	OK = &ErrorX{
		Code:    http.StatusOK,
		Message: "",
	}

	ErrInternal = &ErrorX{
		Code:    http.StatusInternalServerError,
		Reason:  "InternalError",
		Message: "Internal Server Error",
	}
	ErrNotFound = &ErrorX{
		Code:    http.StatusNotFound,
		Reason:  "NotFound",
		Message: "Not Found",
	}
	ErrBind = &ErrorX{
		Code:    http.StatusBadRequest,
		Reason:  "BindError",
		Message: "Bind Request Error",
	}
	ErrInvalidArgument = &ErrorX{
		Code:    http.StatusBadRequest,
		Reason:  "InvalidArgument",
		Message: "Invalid Argument",
	}
	ErrUnauthorized = &ErrorX{
		Code:    http.StatusUnauthorized,
		Reason:  "Unauthorized",
		Message: "Unauthorized",
	}
	ErrPermission = &ErrorX{
		Code:    http.StatusForbidden,
		Reason:  "PermissionDenied",
		Message: "Permission Denied",
	}
	ErrOperation = &ErrorX{
		Code:    http.StatusBadRequest,
		Reason:  "OperationError",
		Message: "Operation Error",
	}
)
