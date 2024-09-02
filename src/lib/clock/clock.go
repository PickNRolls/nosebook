package clock

import (
	"os"
	"time"
)

var testingMockedTime = time.Date(2024, 8, 10, 10, 10, 10, 10, time.UTC)

func Now() time.Time {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "testing" {
		return testingMockedTime
	}

	return time.Now()
}
