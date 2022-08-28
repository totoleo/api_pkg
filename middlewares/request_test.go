package middlewares

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestRecorder(t *testing.T) {
	recorder := RequestRecorder(func() bool {
		return true
	}, func(c *gin.Context) context.Context {

		//like ginex.CacheRPCContext? or ctx.CacheRPCContext😏
		ctx := context.Background()
		for key, _ := range c.Keys {
			v, _ := c.Get(key)
			ctx = context.WithValue(ctx, key, v)
		}
		return ctx
	})

	r := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(r)
	c.Set("t_env", "test-env")
	recorder(c)

}
