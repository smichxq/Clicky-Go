package mtl

import "github.com/kitex-contrib/obs-opentelemetry/provider"

func InitTracing(serviceName string) provider.OtelProvider {
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(serviceName),
		provider.WithExportEndpoint("192.168.3.6:4317"),
		provider.WithInsecure(),
		// change the opentelemetry built-in Metric to prevent conflicts with Prometheus
		provider.WithEnableMetrics(false),
	)

	return p
}
