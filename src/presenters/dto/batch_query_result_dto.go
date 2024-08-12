package dto

import (
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type batchQueryEntry[T any] struct {
	Id     uuid.UUID               `json:"id"`
	Result SingleQueryResultDTO[T] `json:"result"`
}

type BatchQueryResultDTO[T any] struct {
	Err     *errors.Error         `json:"error"`
	Results []*batchQueryEntry[T] `json:"results"`
}
