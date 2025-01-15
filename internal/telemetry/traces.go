package telemetry

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	

)

func InitTraces() {
	// Cria o exportador OTLP
	exp, err := otlptracegrpc.New(context.Background(), otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint("localhost:4317"))
	if err != nil {
		log.Fatalf("failed to create OTLP trace exporter: %v", err)
	}

	// Configura o TracerProvider com o exportador e recursos
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("projeto-otel-na-pratica"), // Substitua pelo nome do seu serviço
		)),
	)

	// Define o TracerProvider global
	otel.SetTracerProvider(tp)

	// Encerra o TracerProvider no final da aplicação
	// (isso deve ser chamado no `main` ou em algum lugar apropriado)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}()
}