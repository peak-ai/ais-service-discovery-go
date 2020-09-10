package redirected

import (
	"errors"

	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
	"github.com/peak-ai/local-plane/config"
)

// QueueAdapter -
type QueueAdapter struct {
	config   config.Config
	resolver Resolver
}

// NewQueueAdapter -
func NewQueueAdapter(config config.Config, resolver Resolver) *QueueAdapter {
	return &QueueAdapter{config, resolver}
}

func (a *QueueAdapter) resolve(service *types.Service) (config.Service, error) {
	srv, ok := a.config[service.Name]
	if !ok {
		return config.Service{}, errors.New("no service reference for address provided, check this service is configured")
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
