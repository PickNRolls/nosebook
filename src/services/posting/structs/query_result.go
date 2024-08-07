package structs

import (
	"nosebook/src/domain/posts"
)

type QueryResult struct {
	Err            error
	RemainingCount int
	Data           []*posts.Post
	Next           string
}
