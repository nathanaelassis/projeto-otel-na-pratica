package main

import (
	"flag"
	"net"
	"net/http"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/config"
	"google.golang.org/grpc"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/cmd/payments"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/cmd/plans"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/cmd/subscriptions"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/cmd/users"

)

func main(){
	configFlag := flag.String("config", "", "path to the config file")
	flag.Parse()

	c, _ := config.LoadConfig(*configFlag)

	mux := http.NewServeMux()

	// starts the gRPC server
	lis, _ := net.Listen("tcp", c.Server.Endpoint.GRPC)
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	{
		payments.Main()
	}

	{
		plans.Main()
	}
	
	{
		subscriptions.Main()
	}

	{
		users.Main()
	}

	go func() {
		_ = grpcServer.Serve(lis)
	}()

	_ = http.ListenAndServe(c.Server.Endpoint.HTTP, mux)

}
