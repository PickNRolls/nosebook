package presenterdto

import (
	"nosebook/src/errors"
)

type FindOut[T any] struct {
	Err        *errors.Error `json:"error,omitempty"`
	Data       []T           `json:"data"`
	TotalCount int           `json:"totalCount"`
	Prev       string        `json:"prev,omitempty"`
	Next       string        `json:"next,omitempty"`
}
