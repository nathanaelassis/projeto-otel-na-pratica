package telemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

func InitTraces() {
	// Configura o exportador OTLP para enviar spans para o endpoint OTLP
	exp, err := otlptracegrpc.New(context.Background())
	if err != nil {
		log.Fatalf("failed to create OTLP trace exporter: %v", err)
	}

	otel.SetTracerProvider(trace.NewTracerProvider(trace.WithBatcher(exp)))

}