package querybuilder

import "github.com/Masterminds/squirrel"

type query = squirrel.SelectBuilder

func New() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

func Union(left query, right query) query {
	sql, args, _ := right.Suffix(")").ToSql()

	return left.
		Prefix("(").
		Suffix(") UNION ("+sql, args...)
}
