package commands

import (
	"nosebook/src/services/commenting/structs"
)

type FindCommentsCommand struct {
	Filter structs.QueryFilter
	Size   *uint
}
