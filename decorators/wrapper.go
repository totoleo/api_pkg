package carry

import (
	"context"
	"errors"
	"net/http"
)

type Response[D any] struct {
	St    int               `json:"st"`
	Msg   string            `json:"msg"`
	Data  D                 `json:"data"`
	Extra map[string]string `json:"extra"` //存放类似 logId, response time 等非业务信息
}

type render interface {
	JSON(int, any)
	Header(key, value string)
}

type Error interface {
	error
	Code() int
	Message() string
}

type Handler[T render] func(c T)
type Endpoint[D any, T render] func(c T) (D, error)

type HandlerV2[T render] func(ctx context.Context, c T)
type EndpointV2[D any, T render] func(ctx context.Context, c T) (D, error)

const SlaHeader = "x-sla"

const SlaError = "0"

// HttpWrapper 适配类似 gin 一类的框架
func HttpWrapper[D any, T render](fn Endpoint[D, T]) Handler[T] {
	return func(c T) {
		data, err := fn(c)
		if err != nil {
			var cErr Error
			var resp Response[D]
			if errors.As(err, &cErr) {
				resp = Response[D]{cErr.Code(), cErr.Message(), data, nil}
			} else {
				resp = Response[D]{500, "系统故障", data, nil}
			}
			c.Header(SlaHeader, SlaError)
			c.JSON(http.StatusOK, resp)
			return
		}
		c.JSON(http.StatusOK, Response[D]{0, "ok", data, nil})
	}
}

// HttpWrapperV2 适配类似 hertz 一类的框架
func HttpWrapperV2[D any, T render](fn EndpointV2[D, T]) HandlerV2[T] {
	return func(ctx context.Context, c T) {
		data, err := fn(ctx, c)
		if err != nil {
			var cErr Error
			var resp Response[D]
			if errors.As(err, &cErr) {
				resp = Response[D]{St: cErr.Code(), Msg: cErr.Message(), Data: data}
			} else {
				resp = Response[D]{St: 500, Msg: "系统故障", Data: data}
			}
			c.Header(SlaHeader, SlaError)
			c.JSON(http.StatusOK, resp)
			return
		}
		c.JSON(http.StatusOK, Response[D]{St: 0, Msg: "ok", Data: data})
	}
}
