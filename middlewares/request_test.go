package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestRecorder(t *testing.T) {
	recorder := RequestRecorder(func() bool {
		return true
	}, func(c *gin.Context) context.Context {
		//like ginex.CacheRPCContext? or ctx.CacheRPCContext😏
		ctx := RequestContext((*ginContext)(c))
		return ctx
	})

	r := httptest.NewRecorder()
	c, e := gin.CreateTestContext(r)
	c.Request, _ = http.NewRequest("POST", "/api", nil)
	c.Request.RequestURI = "/api"
	c.Request.Header = http.Header{}
	c.Request.Header.Set("x-tt-env", "boe_xxx")
	c.Request.Header.Set("x-tt-log", "20220828230501abc223322332233")
	c.Set("t_env", "test-env")
	e.Handlers = append(e.Handlers, recorder)

	recorder(c)

}

type ginContext gin.Context

func (g *ginContext) HandlerName() string {
	return (*gin.Context)(g).HandlerName()
}

func (g *ginContext) GetRequest() *http.Request {
	return g.Request
}
