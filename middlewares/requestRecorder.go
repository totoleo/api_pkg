package middlewares

import (
	"context"
	"log"

	"github.com/bytedance/gopkg/cloud/metainfo"
)

func RequestRecorderMiddleware[R any](enable func() bool) func(ctx context.Context, c R) {
	if !enable() {
		return nil
	}
	return func(ctx context.Context, c R) {
		log.Printf("transient: %+v", metainfo.GetAllValues(ctx))
		log.Printf("persistent: %+v", metainfo.GetAllPersistentValues(ctx))
	}
}

func RequestRecorder[R any](enable func() bool, ctxExtract func(c R) context.Context) func(c R) {
	r := RequestRecorderMiddleware[R](enable)
	if r == nil {
		return nil
	}
	return func(c R) {
		ctx := ctxExtract(c)
		r(ctx, c)
	}
}
