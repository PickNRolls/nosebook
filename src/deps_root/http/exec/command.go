package exec

import (
	"context"
	"nosebook/src/application/services/auth"
	reqcontext "nosebook/src/deps_root/http/req_context"
	"nosebook/src/errors"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func Command[C any, T any](
	serviceMethod func(context.Context, C, *auth.Auth) (T, *errors.Error),
	opts ...func() CommandOption[T],
) func(*gin.Context) {
	return func(ginctx *gin.Context) {
		reqctx := reqcontext.From(ginctx)
		parent := ginctx.Request.Context()

		avoidBinding := false
		var mapper func(out T) any
		var tracer trace.Tracer = noop.Tracer{}
		for _, fn := range opts {
			opt := fn()

			if opt.AvoidBinding() {
				avoidBinding = true
			}

			if opt.Tracer() != nil {
				tracer = opt.Tracer()
			}

			if opt.OutputMapper() != nil {
				mapper = opt.OutputMapper()
			}
		}

		var command C
		if !avoidBinding {
			_, span := tracer.Start(parent, "exec_command.bind_json")
			err := ginctx.ShouldBindJSON(&command)
			span.End()
			if err != nil {
				ginctx.Error(err)
				ginctx.Abort()
				return
			}
		}

		_, span := tracer.Start(parent, "exec_command.service_method")
		out, ok := reqcontext.Handle(serviceMethod(parent, command, reqctx.Auth()))(reqctx)
		span.End()
		reqctx.SetResponseOk(ok)
		if ok {
			var mapped any = out
			if mapper != nil {
				mapped = mapper(out)
			}
			reqctx.SetResponseData(mapped)
		}
	}
}
