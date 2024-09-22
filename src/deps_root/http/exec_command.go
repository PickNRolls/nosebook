package roothttp

import (
	"context"
	"nosebook/src/application/services/auth"
	"nosebook/src/deps_root/http/exec"
	"nosebook/src/errors"

	"github.com/gin-gonic/gin"
)

func execCommand[C any, V any](
  serviceMethod func(context.Context, C, *auth.Auth) (V, *errors.Error),
  root *RootHTTP,
  opts... func() exec.CommandOption[V],
) func(*gin.Context) {
  opts = append(opts, exec.WithTracer[V](root.tracer))
  
	return exec.Command(serviceMethod, opts...) 
}

