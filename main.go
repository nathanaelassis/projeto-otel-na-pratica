package main

import (

	"github.com/dosedetelemetria/projeto-otel-na-pratica/cmd/payments"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/cmd/plans"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/cmd/subscriptions"
	"github.com/dosedetelemetria/projeto-otel-na-pratica/cmd/users"

)

func main(){
	payments.main()
	plans.main()
	subscriptions.main()
	users.main()
}