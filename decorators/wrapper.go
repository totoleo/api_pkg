package decorators

import (
	"context"
	"errors"
	"net/http"
)

// Response 这是一种典型的响应数据结构，当然如果作为一个合格的工具库，这类结构应该由用户制定😂
type Response[D any] struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    D                 `json:"data"`
	Extra   map[string]string `json:"extra"` //存放类似 logId, response time 等非业务信息
}

type render interface {
	JSON(int, any)
	Header(key, value string)
}

// Error 用户可以实现一个满足 Error 接口的结构，携带错误消息和期望展示给用户的信息。如此一来，在业务逻辑只需要判断是否发生错误，以及期望展示给用户什么样的提示
// 至于其它部分，应该交给框架完成。
type Error interface {
	error
	Code() int
	Message() string
}

type Handler[T render] func(c T)
type Endpoint[D any, T render] func(c T) (D, error)

type HandlerV2[T render] func(ctx context.Context, c T)
type EndpointV2[D any, T render] func(ctx context.Context, c T) (D, error)

// SlaHeader x-sla 可以帮助上游的网关来判断接口是否成功完成请求的处理，这样可以在网关完成统一的可用性统计
const SlaHeader = "x-sla"
const slaFailed = "0"

// HttpWrapper 适配类似 gin 一类的框架
func HttpWrapper[D any, T render](fn EndpointV2[D, T], extractor func(c T) context.Context) Handler[T] {
	wp := HttpWrapperV2[D, T](fn)
	return func(c T) {
		ctx := extractor(c)
		wp(ctx, c)
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
				resp = Response[D]{Code: cErr.Code(), Message: cErr.Message(), Data: data}
			} else {
				resp = Response[D]{Code: 500, Message: "系统故障", Data: data}
			}
			c.Header(SlaHeader, slaFailed)
			c.JSON(http.StatusOK, resp)
			return
		}
		c.JSON(http.StatusOK, Response[D]{Code: 0, Message: "ok", Data: data})
	}
}
