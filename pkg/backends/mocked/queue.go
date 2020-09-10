package mocked

import (
	"errors"

	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
)

// QueueAdapter -
type QueueAdapter struct {
	config Config
}

// NewQNewQueueAdapter -
func NewQueueAdapter(config Config) *QueueAdapter {
	return &QueueAdapter{config}
}

func (a *QueueAdapter) resolve(service *types.Service) (Service, error) {

	// In this library, we're passing the signature back in as the name
	srv, ok := a.config[service.Name]
	if !ok {
		return Service{}, errors.New("no service reference for address provided, check this service is configured")
	}

	return srv, nil
}

// QueueWithOpts -
func (a *QueueAdapter) QueueWithOpts(service *types.Service, request types.Request, opts types.Options) (string, error) {
	srv, err := a.resolve(service)
	if err != nil {
		return "", err
	}

	return string(srv.Resolve.MockResponse), nil
}

// ListenWithOpts -
func (a *QueueAdapter) ListenWithOpts(service *types.Service, opts types.Options) (<-chan *types.Response, error) {
	results := make(chan *types.Response)
	srv, err := a.resolve(service)
	if err != nil {
		return nil, err
	}

	go func() {
		results <- &types.Response{
			Body: srv.Resolve.MockResponse,
		}
	}()

	return results, nil
}
