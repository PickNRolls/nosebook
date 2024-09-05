package nullable

import (
	"database/sql"

	"github.com/google/uuid"
)

type nullable[T any] struct {
	Valid bool
	Value T
}

type Bool = nullable[bool]
type Uint64 = nullable[uint64]
type Time = sql.NullTime
type UUID = uuid.NullUUID
