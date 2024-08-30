package presentercomment

import (
	presenterdto "nosebook/src/presenters/dto"
)

type user = presenterdto.User
type likes = presenterdto.Likes
type comment = presenterdto.Comment

type FindByFilterInput struct {
	PostId string
	Next   string
	Prev   string
	Limit  uint64
	Last   bool
}

type FindByFilterOutput = presenterdto.FindOut[*comment]
