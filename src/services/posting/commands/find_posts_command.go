package commands

import (
	"nosebook/src/services/posting/structs"
)

type FindPostsCommand struct {
	Filter structs.QueryFilter
}
