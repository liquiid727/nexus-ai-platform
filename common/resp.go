package common

import "github.com/labstack/echo/v4"

type Resp struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	ErrMsg string      `json:"err_msg"`
	Data   interface{} `json:"data"`
}

func WriteResp(c echo.Context, resp Resp, err error) {
	c.JSON(resp.Code, resp)

}
