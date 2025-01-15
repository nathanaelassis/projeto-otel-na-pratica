// Copyright Dose de Telemetria GmbH
// SPDX-License-Identifier: Apache-2.0

package plans

import (
	"flag"
	"net"
	"net/http"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/app"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/config"
	"google.golang.org/grpc"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/telemetry"
)

func Main() {
	configFlag := flag.String("config", "", "path to the config file")
	flag.Parse()

	// Inicializa a telemetria (traces)
	telemetry.InitTelemetry()

	c, _ := config.LoadConfig(*configFlag)

	// starts the gRPC server
	lis, _ := net.Listen("tcp", c.Server.Endpoint.GRPC)
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	a := app.NewPlan(&c.Plans)
	a.RegisterRoutes(http.DefaultServeMux, grpcServer)

	go func() {
		_ = grpcServer.Serve(lis)
	}()

	_ = http.ListenAndServe(c.Server.Endpoint.HTTP, http.DefaultServeMux)
}
