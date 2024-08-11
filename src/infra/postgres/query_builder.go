package postgres

import "github.com/Masterminds/squirrel"

func NewSquirrel() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
