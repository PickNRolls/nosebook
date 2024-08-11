package generics

import (
	"github.com/google/uuid"
)

type QueryMultiResult[T any] struct {
	Results *struct {
		Id     uuid.UUID
		Result *QuerySingleResult[T]
	}
}
