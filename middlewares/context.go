package middlewares

import (
	"context"
	"net/http"

	"github.com/bytedance/gopkg/cloud/metainfo"
)

type AppContext interface {
	GetRequest() *http.Request
	HandlerName() string
}

// RequestContext 构造应该传递给下游的信息，类似 kitex 的框架会将这类信息放入特定的结构体中，在下游服务中可以使用
func RequestContext[R AppContext](c R) context.Context {

	ctx := context.Background()
	req := c.GetRequest()

	//app info
	ctx = metainfo.WithValue(ctx, "Method", c.HandlerName())
	ctx = metainfo.WithValue(ctx, "Service", CurrentService())

	//http info
	ctx = metainfo.WithPersistentValue(ctx, "request.uri", req.RequestURI)
	ctx = metainfo.WithPersistentValue(ctx, "request.domain", req.URL.Hostname())

	//泳道信息，用以在不同环境中进行路由
	ctx = metainfo.WithPersistentValue(ctx, "k_env", req.Header.Get("x-tt-env"))

	//日志ID
	ctx = metainfo.WithPersistentValue(ctx, "k_log", req.Header.Get("x-tt-log"))

	//save http header info, 透传上游网关认为应该给所有下游的信息
	ctx = metainfo.FromHTTPHeader(ctx, metainfo.HTTPHeader(req.Header))

	return ctx
}

// CurrentService 返回当前容器的服务标识，可以是服务发现的注册名
func CurrentService() string {
	return "fake"
}
