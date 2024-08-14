package generics

import "nosebook/src/errors"

type SingleQueryResult[T any] struct {
	Err            *errors.Error `json:"error"`
	RemainingCount int
	Data           []T `json:"data"`

	Next string `json:"prev"`
	Prev string `json:"next"`
}
