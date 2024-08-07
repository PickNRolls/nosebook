package structs

import (
	"nosebook/src/domain/posts"
	"time"
)

type QueryResult struct {
	Err            error         `json:"error"`
	RemainingCount int           `json:"remainingCount"`
	Data           []*posts.Post `json:"data"`
	Next           time.Time     `json:"next"`
}
