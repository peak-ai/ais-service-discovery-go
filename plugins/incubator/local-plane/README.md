# Local Plane Plugin

Plugin for running services locally with AIS Service Discovery library.

This is a series of mock adapters, which can be used with ais-service-discovery-go. Which allows you to communicate between services locally. Or mock them entirely.

## Use with mocked response

This set-up is the simplest, and allows you to just return some mock responses for specific signatures.
```golang
package main

import (
	"log"

	discovery "github.com/peak-ai/ais-service-discovery-go"
	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
	"github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/backend"
	"github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/config"
)

func main() {

  // Configure a mock endpoint for your service to call
	config := config.Config{
		"test.service->queue": config.Service{
			Type: "queue",
			Resolve: config.Resolver{
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

```


## Use with a redirected endpoint

This set-up uses the 'redirect' backend, which allows you to define an internal type and endpoint to call instead
of actual AWS services. In the example below, we set the resolver type to 'http-post', and define an endpoint for 
a specific service address.

```golang
package main

import (
        "log"
        discovery "github.com/peak-ai/ais-service-discovery-go"
        "github.com/peak-ai/ais-service-discovery-go/pkg/types"
        "github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/backend"
        "github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/config"
)

func main() {

  // Configure a mock endpoint for your service to call
	config := config.Config{
		"test.service->queue": config.Service{
			Type: "queue",
			Resolve: config.Resolver{
				Type:         "http-post",
				Endpoint:     "http://localhost:8080",
			},
		},
	}

  d := discovery.NewDiscovery(discovery.WithAWSBackend())

  // If env is local, use the mocked back-end instead
  if os.Getenv("ENV") == "local" {
    d := discovery.NewDiscovery(backend.WithRedirectedBackend(config))
  }

  id, err := d.Queue("test.service->queue", types.Request{})
  if err != nil {
    log.Panic(err)
  }

  log.Println("Response: ", id)
}

```

## TODO
- [ ] - Create adapters for events
- [ ] - Better tests
- [ ] - Error simulation scenarios

Other ideas:
- local config file perhaps, json or hcl, such as .local-plane, or .service-discovery, which contains some environment based config.
