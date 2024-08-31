package cursor

import (
	"fmt"
	"nosebook/src/errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

func Encode(timestamp time.Time, id uuid.UUID) string {
	return fmt.Sprintf("%v/%v", id, timestamp.Format(time.RFC3339Nano))
}

func Decode(str string) (time.Time, uuid.UUID, *errors.Error) {
	substrings := strings.Split(str, "/")
	idStr := substrings[0]
	timeStr := substrings[1]

	timestamp, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return time.Time{}, uuid.Nil, errors.New("Cursor Decode Error", "Invalid cursor")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return time.Time{}, uuid.Nil, errors.New("Cursor Decode Error", "Invalid cursor")
	}

	return timestamp, id, nil
}
