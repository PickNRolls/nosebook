package clock

import (
	"nosebook/src/lib/config"
	"time"
)

var testingMockedTime = time.Date(2024, 8, 10, 10, 10, 10, 0, time.UTC)

func Now() time.Time {
	if config.Env.IsTesting() {
		return testingMockedTime
	}

	return time.Now()
}
