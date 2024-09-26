package config

import "os"

type config struct {
	Env     appEnv
	Tracing tracing
}

var Config = config{
	Env: appEnv{env: os.Getenv("APP_ENV")},
	Tracing: tracing{enabled: func() bool {
		variable := os.Getenv("TRACING_ENABLED")
		return variable != "" && variable != "0"
	}()},
}
