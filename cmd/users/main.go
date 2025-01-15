package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/app"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/config"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	// Lê o caminho do arquivo de configuração
	configFlag := flag.String("config", "", "path to the config file")
	flag.Parse()

	// Inicializa a telemetria (traces)
	telemetry.InitTelemetry()

	// Carrega a configuração
	c, err := config.LoadConfig(*configFlag)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Cria a aplicação e registra as rotas
	a := app.NewUser(&c.Users)
	a.RegisterRoutes(http.DefaultServeMux)

	// Instrumenta o servidor HTTP com o middleware do OpenTelemetry
	otelHandler := otelhttp.NewHandler(http.DefaultServeMux, "HTTP Server")

	// Inicia o servidor HTTP
	log.Printf("Starting server at %s", c.Server.Endpoint.HTTP)
	if err := http.ListenAndServe(c.Server.Endpoint.HTTP, otelHandler); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}