package main

import (
	"flag"
	"juno/pkg/api/client"
	"juno/pkg/ranag/router"
	"juno/pkg/ranag/service"
	"time"

	"juno/pkg/ranag/handler"

	"github.com/sirupsen/logrus"
)

func main() {

	var apiURL string
	var port string
	flag.StringVar(&apiURL, "api-url", "http://127.0.0.1:8080", "API URL")
	flag.StringVar(&port, "port", "6060", "Port to run the server on")
	flag.Parse()

	if apiURL == "" {
		panic("api-url is required")
	}

	s := service.New(
		service.WithApiClient(
			client.New(
				apiURL,
			),
		),

		service.WithLogger(
			logrus.New(),
		),

		service.WithShardFetchInterval(time.Minute),
	)

	h := handler.New(s)

	r := router.New(h)

	r.Run(":" + port)
}
