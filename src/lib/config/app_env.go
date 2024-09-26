package config

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
