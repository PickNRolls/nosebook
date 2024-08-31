package roothttp

import (
	reqcontext "nosebook/src/deps_root/http/req_context"
	"nosebook/src/errors"
)

func handle[T any](data T, err *errors.Error) func(*reqcontext.ReqContext) (T, bool) {
	return reqcontext.Handle(data, err)
}
