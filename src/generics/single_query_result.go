package generics

import "nosebook/src/errors"

type SingleQueryResult[T any] struct {
	Err            *errors.Error
	RemainingCount int
	Data           []T

	Next string
	Prev string
}
