package main

import (
	"log"

	discovery "github.com/peak-ai/ais-service-discovery-go"
	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
	"github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/backend"
	"github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/config"
)

func main() {
	config := config.Config{
		"namespace.service->queue": config.Service{
			Type: "queue",
			Resolve: config.Resolver{
				MockResponse: []byte("queue response"),
			},
		},
		"namespace.service->request": config.Service{
			Type: "function",
			Resolve: config.Resolver{
				MockResponse: []byte("request response"),
			},
		},
	}

	sd := discovery.NewDiscovery(backend.WithMockedBackend(config))
	queueResponse, err := sd.Queue("namespace.service->queue", types.Request{})
	if err != nil {
		log.Panic(err)
	}

	log.Println(queueResponse)

	requestResponse, err := sd.Request("namespace.service->request", types.Request{})
	if err != nil {
		log.Panic(err)
	}

	log.Println(requestResponse)
}
