package presenterdto

import (
	"nosebook/src/errors"
)

type FindOut[T any] struct {
	Err  *errors.Error `json:"error,omitempty"`
	Data []T           `json:"data"`
	Next string        `json:"next,omitempty"`
}
