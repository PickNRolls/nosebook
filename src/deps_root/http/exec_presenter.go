package roothttp

import (
	"context"
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/application/services/auth"
	reqcontext "nosebook/src/deps_root/http/req_context"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type MapBuilder[T any] interface {
	BuildFromMap(m map[string]any) T
}

type declarationType = string

const (
	STRING declarationType = "STRING"
	UINT64 declarationType = "UINT64"
)

type presenterOption struct {
	Type     declarationType
	Required bool
}

func execDefaultPresenter[In any, Out any](
	fn func(context.Context, In, *auth.Auth) *presenterdto.FindOut[Out],
	declarations map[string]presenterOption,
	builder MapBuilder[In],
	tracer trace.Tracer,
) func(*gin.Context) {
	return func(ctx *gin.Context) {
		c, span := tracer.Start(ctx.Request.Context(), "presenter_handler")
		defer span.End()

		reqctx := reqcontext.From(ctx)

		m := make(map[string]any)
		for key, decl := range declarations {
			if decl.Type == STRING {
				m[key] = ctx.Query(key)
			}

			if decl.Type == UINT64 {
				if !decl.Required {
					value, ok := reqctx.QueryNullableUint64(key)
					if !ok {
						return
					}
					m[key] = value
				} else {
					value, ok := reqctx.QueryUint64(key)
					if !ok {
						return
					}
					m[key] = value
				}
			}
		}

		input := builder.BuildFromMap(m)
		output := fn(c, input, reqctx.Auth())

		_, ok := handle(output, output.Err)(reqctx)
		if !ok {
			return
		}

		reqctx.SetResponseData(output)
		reqctx.SetResponseOk(true)
	}
}
