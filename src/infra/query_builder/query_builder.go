package querybuilder

import (
	"github.com/Masterminds/squirrel"
)

type query = squirrel.SelectBuilder

type Opt interface {
  OmitPlaceholder() bool
}

type OptFn func() Opt

type omitPlaceholderOpt struct {
}

func (this *omitPlaceholderOpt) OmitPlaceholder() bool {
  return true
}

func OmitPlaceholder() Opt {
  return &omitPlaceholderOpt{}
}

func New(fns ...OptFn) squirrel.StatementBuilderType {
	out := squirrel.StatementBuilder

  omitPlaceholder := false

  for _, fn := range fns {
    opt := fn()
    if opt.OmitPlaceholder() {
      omitPlaceholder = true
    }
  }

	if !omitPlaceholder {
		out = out.PlaceholderFormat(squirrel.Dollar)
	}

	return out
}

func Union(left query, right query) query {
	sql, args, _ := right.Suffix(")").ToSql()

	return left.
		Prefix("(").
		Suffix(") UNION ("+sql, args...)
}

