package config

import "os"

var Env = &appEnv{env: os.Getenv("APP_ENV")}

var Tracing = &tracing{
	enabled: func() bool {
		variable := os.Getenv("TRACING_ENABLED")
		return variable != "" && variable != "0"
	}(),

	exporter: func() string {
		exporter := os.Getenv("TRACING_EXPORTER")
		if exporter != "jaeger" && exporter != "stdout" {
			exporter = "stdout"
		}
		return exporter
	}(),
}

var DBName = os.Getenv("POSTGRES_DB")
