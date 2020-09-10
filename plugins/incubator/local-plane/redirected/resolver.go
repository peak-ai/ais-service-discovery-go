//go:generate mockery --name=Resolver --case=underscore --output=mocks
package redirected

import (
	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
	"github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/config"
)

// Resolver is an interface for a type which takes
// a resolver config and returns a response
type Resolver interface {
	Resolve(resolver config.Resolver) (*types.Response, error)
}
