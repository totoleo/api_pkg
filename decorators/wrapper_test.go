package decorators_test

import (
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/cloudwego/hertz"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/gin-gonic/gin"

	"github.com/totoleo/api_pkg/decorators"
)

func TestHttpWrapper(t *testing.T) {
	gin.SetMode("release")
	wp := decorators.HttpWrapper(ginHandler)
	writer := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(writer)

	wp(c)

	result := writer.Result()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(body), len(body))
	t.Log(writer.Header())
}

func TestHertz(t *testing.T) {

	wp := decorators.HttpWrapperV2(hertzHandler)

	c := app.NewContext(1)

	wp(context.Background(), c)

	body := c.Response.Body()
	t.Log(string(body), len(body))
	t.Log(string(c.Response.Header.Header()))
}

func ginHandler(c *gin.Context) (*Vo, error) {
	return &Vo{
		Id:   1,
		Name: "gin@" + gin.Version,
	}, errors.New("gin")
}

func hertzHandler(ctx context.Context, c *app.RequestContext) (Vo, error) {
	return Vo{
		Id:   2,
		Name: "hertz@" + hertz.Version,
	}, errors.New("hertz")
}

type Vo struct {
	Id   int64 `json:"id"`
	Name string
}
