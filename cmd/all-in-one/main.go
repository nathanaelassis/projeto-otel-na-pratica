package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/app"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/config"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
)

func main() {
	configFlag := flag.String("config", "", "path to the config file")
	flag.Parse()

	// Inicializa a telemetria (traces)
	telemetry.InitTelemetry()

	c, err := config.LoadConfig(*configFlag)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	mux := http.NewServeMux()

	// starts the gRPC server
	lis, err := net.Listen("tcp", c.Server.Endpoint.GRPC)
	if err != nil {
		log.Fatalf("failed to start gRPC listener: %v", err)
	}

	// Adiciona interceptores para instrumentar gRPC
	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	opts = append(opts, grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()))
	grpcServer := grpc.NewServer(opts...)

	{
		a := app.NewUser(&c.Users)
		// Instrumenta as rotas HTTP com otelhttp
		a.RegisterRoutes(mux)
	}
	{
		a := app.NewPlan(&c.Plans)
		a.RegisterRoutes(mux, grpcServer)
	}
	{
		a, err := app.NewPayment(&c.Payments)
		if err != nil {
			panic(err)
		}
		a.RegisterRoutes(mux)
		defer func() {
			_ = a.Shutdown()
		}()
	}
	{
		a := app.NewSubscription(&c.Subscriptions)
		a.RegisterRoutes(mux)
	}

	// Inicia o servidor gRPC em uma goroutine
	go func() {
		log.Printf("Starting gRPC server at %s", c.Server.Endpoint.GRPC)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to start gRPC server: %v", err)
		}
	}()

	// Instrumenta o servidor HTTP com o middleware do OpenTelemetry
	otelHandler := otelhttp.NewHandler(mux, "HTTP Server")

	// Inicia o servidor HTTP
	log.Printf("Starting HTTP server at %s", c.Server.Endpoint.HTTP)
	if err := http.ListenAndServe(c.Server.Endpoint.HTTP, otelHandler); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}