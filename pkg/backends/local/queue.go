package local

import (
	"errors"

	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
)

// QueueAdapter -
type QueueAdapter struct {
	config   Config
	resolver *ResolverService
}

// NewQueueAdapter -
func NewQueueAdapter(config Config, resolver *ResolverService) *QueueAdapter {
	return &QueueAdapter{config, resolver}
}

func (a *QueueAdapter) resolve(service *types.Service) (Service, error) {
	srv, ok := a.config[service.Addr]
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

	response, err := a.resolver.Resolve(srv.Resolve)

	return string(response.Body), err
}

// ListenWithOpts -
func (a *QueueAdapter) ListenWithOpts(service *types.Service, opts types.Options) (<-chan *types.Response, error) {
	results := make(chan *types.Response)
	srv, err := a.resolve(service)
	if err != nil {
		return nil, err
	}

	go func() {
		response, _ := a.resolver.Resolve(srv.Resolve)
		// @todo - need to handle error here
		results <- response
	}()

	return results, nil
}
