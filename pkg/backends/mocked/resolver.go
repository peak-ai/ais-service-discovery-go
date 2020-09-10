package mocked

import (
	"net/http"

	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
)

// NewResolver -
func NewResolver() *ResolverService {
	return &ResolverService{}
}

// Resolver -
type ResolverService struct{}

// Resolve -
func (r *ResolverService) Resolve(resolver Resolver) (*types.Response, error) {
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
