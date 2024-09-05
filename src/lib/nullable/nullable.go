package nullable

import (
	"database/sql"
)

type nullable[T any] struct {
	Valid bool
	Value T
}

type Bool = nullable[bool]
type Uint64 = nullable[uint64]
type Time = sql.NullTime
