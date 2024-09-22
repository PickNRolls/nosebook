package commandresult

import (
	"nosebook/src/errors"
)

type Result struct {
	Ok     bool            `json:"ok"`
	Errors []*errors.Error `json:"errors,omitempty"`
	Data   any             `json:"data,omitempty"`
}

