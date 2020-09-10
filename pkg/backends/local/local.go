package local

import (
	discovery "github.com/peak-ai/ais-service-discovery-go"
)

// Resolver is the action to perform instead of the
// original action. I.e http post request
type Resolver struct {
	Type         string `json:"type"`
	Endpoint     string `json:"endpoint"`
	Payload      []byte `json:"payload"`
	MockResponse []byte `json:"mock_response"`
}

// Service type is the original type,
// i.e. queue, event etc.
// Resolver is the action to perform instead.
type Service struct {
	Type    string   `json:"type"`
	Resolve Resolver `json:"resolve"`
}

// Config the key is a service address, and the value
// is the metadata used to figure out what to do when
// that service is called.
type Config map[string]Service

func Factory(config Config) discovery.Option {
	resolver := NewResolver()
	return func(args *discovery.Options) {
		args.Locator = NewLocator(config)
		args.QueueAdapter = NewQueueAdapter(config, resolver)
	}
}
