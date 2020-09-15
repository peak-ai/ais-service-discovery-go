package mocked

import (
	"errors"

	"github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/config"
	"github.com/peak-ai/ais-service-discovery-go/types"
)

type FunctionAdapter struct {
	config config.Config
}

// NewFunctionAdapter -
func NewFunctionAdapter(config config.Config) *FunctionAdapter {
	return &FunctionAdapter{config}
}

func (f *FunctionAdapter) resolve(service *types.Service) (config.Service, error) {
	srv, ok := f.config[service.Name]
	if !ok {
		return config.Service{}, errors.New("no function found for that address")
	}

	return srv, nil
}

// CallWithOpts mocks a cloud function call
func (f *FunctionAdapter) CallWithOpts(service *types.Service, request types.Request, opts types.Options) (*types.Response, error) {
	srv, err := f.resolve(service)
	if err != nil {
		return nil, err
	}

	return &types.Response{
		Body: srv.Resolve.MockResponse,
	}, nil
}

// Call a cloud function
func (f *FunctionAdapter) Call(service *types.Service, request types.Request) (*types.Response, error) {
	return f.CallWithOpts(service, request, types.Options{})
}
