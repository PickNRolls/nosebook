package config

type tracing struct {
	enabled bool
}

func (this *tracing) Enabled() bool {
	return this.enabled
}
