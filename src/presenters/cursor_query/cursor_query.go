package cursorquery

import (
	"nosebook/src/boolean"
	"nosebook/src/errors"
	"nosebook/src/infra/postgres"
	"nosebook/src/presenters/cursor"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type destType interface {
	Timestamp() time.Time
	ID() uuid.UUID
}

type Input struct {
	Query    squirrel.SelectBuilder
	Limit    uint64
	Prev     string
	Next     string
	Last     bool
	OrderAsc bool
}

type Cursors struct {
	Next string
	Prev string
}

func newError(message string) *errors.Error {
	return errors.New("Cursor Query Error", message)
}

func errorFrom(err error) *errors.Error {
	return newError(err.Error())
}

func addOrder(query squirrel.SelectBuilder, asc bool) squirrel.SelectBuilder {
	if asc {
		return query.OrderBy("created_at asc, id asc")
	}

	return query.OrderBy("created_at desc, id desc")
}

func Do[T destType](db *sqlx.DB, input *Input, dest *[]T) (*Cursors, *errors.Error) {
	if input.Prev != "" && input.Next != "" {
		return nil, newError("Использовать Prev и Next одновременно невозможно")
	}

	next := input.Next
	prev := input.Prev

	query := input.Query.Limit(input.Limit + 1)
	query = addOrder(query, boolean.Xor(input.OrderAsc, input.Last))

	if input.Last {
		next = ""
		prev = ""

		qb := postgres.NewSquirrel()
		query := qb.Select("*").
			FromSelect(query, "inner")
		query = addOrder(query, input.OrderAsc)
	}

	if next != "" {
		timestamp, id, err := cursor.Decode(next)
		if err != nil {
			return nil, errorFrom(err)
		}

		if input.OrderAsc {
			query = query.Where("(created_at, id) > (?, ?)", timestamp, id)
		} else {
			query = query.Where("(created_at, id) < (?, ?)", timestamp, id)
		}
	}

	if prev != "" {
		timestamp, id, err := cursor.Decode(prev)
		if err != nil {
			return nil, errorFrom(err)
		}

		if input.OrderAsc {
			query = query.Where("(created_at, id) < (?, ?)", timestamp, id)
		} else {
			query = query.Where("(created_at, id) > (?, ?)", timestamp, id)
		}
	}

	sql, args, _ := query.ToSql()
	err := db.Select(dest, sql, args...)
	if err != nil {
		return nil, errorFrom(err)
	}

	slice := *dest
	lengthGreaterLimit := len(slice) > int(input.Limit)
	shouldBeNext := lengthGreaterLimit
	shouldBePrev := lengthGreaterLimit && (prev != "" || input.Last)

	if lengthGreaterLimit {
		*dest = slice[:len(slice)-1]
		slice = *dest
	}

	output := &Cursors{}

	if shouldBeNext {
		last := slice[len(slice)-1]
		output.Next = cursor.Encode(last.Timestamp(), last.ID())
	}

	if shouldBePrev {
		first := slice[0]
		output.Prev = cursor.Encode(first.Timestamp(), first.ID())
	}

	return output, nil
}
