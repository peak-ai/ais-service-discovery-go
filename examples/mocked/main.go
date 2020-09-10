package main

import (
	discovery "github.com/peak-ai/ais-service-discovery-go"
	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
	"github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/backend"
	"github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/config"
)

func main() {
	config := config.Config{}
	sd := discovery.NewDiscovery(backend.WithMockedBackend(config))
	sd.Queue("namespace.service->test", types.Request{})
}
