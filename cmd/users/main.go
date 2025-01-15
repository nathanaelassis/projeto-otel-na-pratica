// Copyright Dose de Telemetria GmbH
// SPDX-License-Identifier: Apache-2.0

package users

import (
	"flag"
	"net/http"

	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/app"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/config"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/internal/telemetry"
)

func Main() {
	configFlag := flag.String("config", "", "path to the config file")
	flag.Parse()

	// Inicializa a telemetria (traces)
	telemetry.InitTelemetry()

	c, _ := config.LoadConfig(*configFlag)

	a := app.NewUser(&c.Users)
	a.RegisterRoutes(http.DefaultServeMux)
	_ = http.ListenAndServe(c.Server.Endpoint.HTTP, http.DefaultServeMux)
}


