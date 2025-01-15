// Copyright Dose de Telemetria GmbH
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"net/http"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/api"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/config"
	grpchandler "github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/handler/grpc"
	planhttp "github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/handler/http"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/store"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/pkg/store/memory"
	"google.golang.org/grpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp" //Importar a dependencia otelhttp
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	
)

type Plan struct {
	Handler     *planhttp.PlanHandler
	GRPCHandler api.PlanServiceServer
	Store       store.Plan
}


func NewPlan(*config.Plans) *Plan {
	store := memory.NewPlanStore()
	return &Plan{
		Handler:     planhttp.NewPlanHandler(store),
		GRPCHandler: grpchandler.NewPlanServer(store),
		Store:       store,
	}
}

func (a *Plan) RegisterRoutes(mux *http.ServeMux, grpcSrv *grpc.Server) {
	// Instrumenta as rotas HTTP com o middleware otelhttp
	mux.Handle("GET /plans", otelhttp.NewHandler(http.HandlerFunc(a.Handler.List), "ListPlans"))
	mux.Handle("POST /plans", otelhttp.NewHandler(http.HandlerFunc(a.Handler.Create), "CreatePlan"))
	mux.Handle("GET /plans/{id}", otelhttp.NewHandler(http.HandlerFunc(a.Handler.Get), "GetPlan"))
	mux.Handle("PUT /plans/{id}", otelhttp.NewHandler(http.HandlerFunc(a.Handler.Update), "UpdatePlan"))
	mux.Handle("DELETE /plans/{id}", otelhttp.NewHandler(http.HandlerFunc(a.Handler.Delete), "DeletePlan"))

	// Instrumenta o servidor gRPC com interceptores do OpenTelemetry
	grpcSrv = grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	api.RegisterPlanServiceServer(grpcSrv, a.GRPCHandler)
}