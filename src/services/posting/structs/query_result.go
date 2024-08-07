package structs

import (
	"nosebook/src/domain/posts"
)

type QueryResult struct {
	Err            error         `json:"error"`
	RemainingCount int           `json:"remainingCount"`
	Data           []*posts.Post `json:"data"`
	Next           string        `json:"next"`
}
