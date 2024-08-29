package commandresult

import (
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type Result struct {
	Ok     bool            `json:"ok"`
	Errors []*errors.Error `json:"errors,omitempty"`
	Data   any             `json:"data,omitempty"`
}

func FromCommandReturn(data any, err *errors.Error) *Result {
	if err != nil {
		return Fail(err)
	}

	output := Ok()
	if id, ok := data.(uuid.UUID); ok {
		output.WithId(id)
		return output
	}

	return output.WithData(data)
}

func Ok() *Result {
	return &Result{
		Ok: true,
	}
}

func Fail(err *errors.Error) *Result {
	return &Result{
		Errors: []*errors.Error{err},
		Ok:     false,
	}
}

func (this *Result) WithData(data any) *Result {
	this.Data = data
	return this
}

func (this *Result) WithId(id uuid.UUID) *Result {
	this.Data = struct {
		Id uuid.UUID `json:"id"`
	}{
		Id: id,
	}
	return this
}
