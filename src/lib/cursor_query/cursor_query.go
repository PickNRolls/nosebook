package cursorquery

import (
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	"nosebook/src/lib/boolean"
	"nosebook/src/lib/cursor"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Order[T any] interface {
	Column() string
	Timestamp(dest T) time.Time
	Id(dest T) uuid.UUID
	Asc() bool
}

type Input[T any] struct {
	Query squirrel.SelectBuilder
	Limit uint64
	Prev  string
	Next  string
	Last  bool
	Order Order[T]
}

type Output struct {
	TotalCount int
	Next       string
	Prev       string
}

const DEFAULT_LIMIT = 10
const MAX_LIMIT = 20

func newError(message string) *errors.Error {
	return errors.New("Cursor Query Error", message)
}

func errorFrom(err error) *errors.Error {
	return newError(err.Error())
}

func addOrder(query squirrel.SelectBuilder, column string, asc bool) squirrel.SelectBuilder {
	if asc {
		return query.OrderBy(column + " asc, id asc")
	}

	return query.OrderBy(column + " desc, id desc")
}

func Do[T any](db *sqlx.DB, input *Input[T], dest *[]T) (*Output, *errors.Error) {
	if input.Prev != "" && input.Next != "" {
		return nil, newError("Использовать Prev и Next одновременно невозможно")
	}

	if input.Limit > MAX_LIMIT {
		return nil, newError("Максимальный Limit = " + string(MAX_LIMIT))
	}

	order := input.Order
	orderColumn := order.Column()
	if orderColumn == "" {
		orderColumn = "created_at"
	}

	next := input.Next
	prev := input.Prev
	limit := input.Limit
	if limit == 0 {
		limit = DEFAULT_LIMIT
	}

	qb := querybuilder.New()
	query := qb.Select("*").
		FromSelect(input.Query, "sub").
		Limit(limit + 1)

	query = addOrder(query, orderColumn, boolean.Xor(order.Asc(), input.Last))

	if input.Last {
		next = ""
		prev = ""

		qb := querybuilder.New()
		query := qb.Select("*").
			FromSelect(query, "inner")
		query = addOrder(query, orderColumn, order.Asc())
	}

	if next != "" {
		timestamp, id, err := cursor.Decode(next)
		if err != nil {
			return nil, errorFrom(err)
		}

		if order.Asc() {
			query = query.Where("("+orderColumn+", id) > (?, ?)", timestamp, id)
		} else {
			query = query.Where("("+orderColumn+", id) < (?, ?)", timestamp, id)
		}
	}

	if prev != "" {
		timestamp, id, err := cursor.Decode(prev)
		if err != nil {
			return nil, errorFrom(err)
		}

		if order.Asc() {
			query = query.Where("("+orderColumn+", id) < (?, ?)", timestamp, id)
		} else {
			query = query.Where("("+orderColumn+", id) > (?, ?)", timestamp, id)
		}
	}

	sql, args, _ := query.ToSql()
	err := db.Select(dest, sql, args...)
	if err != nil {
		return nil, errorFrom(err)
	}

	totalCount, err := func() (int, *errors.Error) {
		count := struct {
			Count int `db:"count"`
		}{}
		query := qb.Select("count(*)").FromSelect(input.Query, "sub")
		sql, args, _ := query.ToSql()
		err := db.Get(&count, sql, args...)
		if err != nil {
			return 0, errorFrom(err)
		}

		return count.Count, nil
	}()

	slice := *dest
	lengthGreaterLimit := len(slice) > int(limit)
	shouldBeNext := lengthGreaterLimit
	shouldBePrev := lengthGreaterLimit && (prev != "" || input.Last)

	if lengthGreaterLimit {
		*dest = slice[:len(slice)-1]
		slice = *dest
	}

	output := &Output{
		TotalCount: totalCount,
	}

	if shouldBeNext {
		last := slice[len(slice)-1]
		output.Next = cursor.Encode(input.Order.Timestamp(last), input.Order.Id(last))
	}

	if shouldBePrev {
		first := slice[0]
		output.Prev = cursor.Encode(input.Order.Timestamp(first), input.Order.Id(first))
	}

	return output, nil
}
