package mocked

import (
	"errors"
	"fmt"

	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
	"github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/config"
)

// NewLocator -
func NewLocator(config config.Config) *Locator {
	return &Locator{config}
}

// Locator -
type Locator struct {
	config config.Config
}

// Discover -
func (l *Locator) Discover(signature *types.Signature) (*types.Service, error) {
	// This is kind weird, we have to rebuild the signature again here,
	// because at this point it's parsed typically.
	s := fmt.Sprintf("%s.%s->%s", signature.Namespace, signature.Service, signature.Instance)
	srv, ok := l.config[s]
	if !ok {
		return nil, errors.New("error finding service in config")
	}

	return &types.Service{
		Name: s, // This is kinda naughty, it isn't technically the 'name' it's the signature
		Addr: srv.Resolve.Endpoint,
	}, nil
}
