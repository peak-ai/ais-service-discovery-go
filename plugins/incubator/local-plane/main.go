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
		"test.service->queue": config.Service{
			Type: "queue",
			Resolve: config.Resolver{
				Type:         "http-post",
				Endpoint:     "http://localhost:8080",
				MockResponse: []byte("testing 123"),
			},
		},
	}
	d := discovery.NewDiscovery(backend.WithMockedBackend(config))

	id, err := d.Queue("test.service->queue", types.Request{})
	if err != nil {
		log.Panic(err)
	}

	log.Println("id: ", id)
}
