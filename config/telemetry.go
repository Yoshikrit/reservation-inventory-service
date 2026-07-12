package config

type TelemetryConfig struct {
	OtelEndpoint    string `env:"OTEL_ENDPOINT"     envDefault:"tempo:4317"`
	OtelServiceName string `env:"OTEL_SERVICE_NAME" envDefault:"inventory"`
}
