package generics

import "nosebook/src/errors"

type QueryResult[T any] struct {
	Err            *errors.Error
	RemainingCount int
	Data           []T

	Next string
	Prev string
}
