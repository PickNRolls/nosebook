package presenterdto

import (
	"encoding/json"
	"nosebook/src/errors"
)

type Event struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

func (this *Event) ToJson() ([]byte, *errors.Error) {
	json, err := errors.Using(json.Marshal(this))
	if err != nil {
		return nil, err
	}
	return json, err
}
