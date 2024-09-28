package config

type tracing struct {
	enabled bool
  exporter string
}

func (this *tracing) Enabled() bool {
	return this.enabled
}

func (this *tracing) IsJaegerExporter() bool {
  return this.exporter == "jaeger"
}

