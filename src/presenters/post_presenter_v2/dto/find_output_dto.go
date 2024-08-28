package dto

import (
	"nosebook/src/errors"
)

type FindOutputDTO struct {
	Err  *errors.Error `json:"error,omitempty"`
	Data []*PostDTO    `json:"data,omitempty"`
	Next string        `json:"next,omitempty"`
}
