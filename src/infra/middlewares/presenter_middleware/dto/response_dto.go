package dto

import "nosebook/src/errors"

type ResponseDTO struct {
	Errors []*errors.Error `json:"errors"`
	Data   any             `json:"data"`
}
