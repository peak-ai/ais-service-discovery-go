package resolvers

import (
	"net/http"

	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
	"github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/config"
)

// NewResolver -
func NewResolver() *Resolver {
	return &Resolver{}
}

// Resolver -
type Resolver struct{}

// Resolve -
func (r *Resolver) Resolve(resolver config.Resolver) (*types.Response, error) {
	switch resolver.Type {
	case "http-post":
		client := &http.Client{}
		r := NewHttpResolver(client)
		return r.Post(resolver)

	case "http-get":
		client := &http.Client{}
		r := NewHttpResolver(client)
		return r.Get(resolver)

	default:
		client := &http.Client{}
		r := NewHttpResolver(client)
		return r.Get(resolver)
	}
}
