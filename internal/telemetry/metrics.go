package telemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	//"go.opentelemetry.io/otel/sdk/resource"
	//"go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitMetrics() {
	// Configura o exportador OTLP para métricas
	exp, err := otlpmetricgrpc.New(context.Background())
	if err != nil {
		log.Fatalf("failed to create OTLP trace exporter: %v", err)
	}

	// Configura o PeriodicReader para coletar e exportar métricas periodicamente
	reader := metric.NewPeriodicReader(exp)

	otel.SetMeterProvider(metric.NewMeterProvider(metric.WithReader(reader)))

}