package common

import "fmt"

type Errno struct {
	Code int
	Msg  string
	Err  error
}

// usage
// err := common.Errno{Code: 400, Msg: "bad request"}
// err = err.WithErr(errors.New("invalid parameter"))
// fmt.Println(err.Error()) // Output: bad request: invalid parameter

func (err Errno) Error() string {
	if err.Err != nil {
		return fmt.Sprintf("%s: %v", err.Msg, err.Err)
	}

	return err.Msg
}
func (e Errno) Unwrap() error {
	return e.Err
}

func (err Errno) WithErr(rawErr error) Errno {
	err.Err = rawErr
	return err
}

var (
	OK              = Errno{Code: 200, Msg: "ok"}
	ServerError     = Errno{Code: 500, Msg: "server error"}
	BadRequest      = Errno{Code: 400, Msg: "bad request"}
	AuthError       = Errno{Code: 401, Msg: "auth error"}
	PermissionError = Errno{Code: 403, Msg: "permission error"}
	DatabaseError   = Errno{Code: 10000, Msg: "database error"}
	RedisError      = Errno{Code: 10001, Msg: "redis error"}

	UserNotFoundErr   = Errno{Code: 10002, Msg: "user not found"}
	InvalidCaptchaErr = Errno{Code: 10003, Msg: "invalid captcha"}
)
