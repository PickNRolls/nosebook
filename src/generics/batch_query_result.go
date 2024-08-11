package generics

import (
	"github.com/google/uuid"
)

type BatchQueryResult[T any] struct {
	Results *struct {
		Id     uuid.UUID
		Result *SingleQueryResult[T]
	}
}
