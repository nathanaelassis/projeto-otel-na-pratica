// Copyright Dose de Telemetria GmbH
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"net"
	"net/http"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/app"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/config"
	"google.golang.org/grpc"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/telemetry"
)


func main() {
	configFlag := flag.String("config", "", "path to the config file")
	flag.Parse()

	c, _ := config.LoadConfig(*configFlag)

	mux := http.NewServeMux()

	// starts the gRPC server
	lis, _ := net.Listen("tcp", c.Server.Endpoint.GRPC)
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// Inicializa a telemetria (traces)
	

	{
		telemetry.InitTelemetry()
		a := app.NewUser(&c.Users)
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

	go func() {
		_ = grpcServer.Serve(lis)
	}()

	_ = http.ListenAndServe(c.Server.Endpoint.HTTP, mux)
}
