package dto

import "nosebook/src/errors"

type QuerySingleResultDTO[T any] struct {
	Err            *errors.Error `json:"error"`
	RemainingCount int           `json:"remainingCount"`
	Data           []T           `json:"data"`
	Next           string        `json:"next"`
}
