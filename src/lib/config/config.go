package config

import "os"

type appEnv struct {
	env string
}

func (this *appEnv) IsProduction() bool {
	return this.env == "prod"
}
func (this *appEnv) IsTesting() bool {
	return this.env == "testing"
}
func (this *appEnv) IsDevelopment() bool {
  return this.env == "dev"
}

type config struct {
	Env appEnv
}

var Config = config{
	Env: appEnv{env: os.Getenv("APP_ENV")},
}

