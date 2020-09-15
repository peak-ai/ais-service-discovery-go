package redirected

import (
	"testing"

	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
	"github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/config"
	"github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/redirected/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestCanQueue -
func TestCanQueue(t *testing.T) {
	response := &types.Response{
		Body: []byte("test"),
	}
	resolver := &mocks.Resolver{}
	resolver.On("Resolve", mock.AnythingOfType("config.Resolver")).
		Return(response, nil)

	c := config.Config{
		"latest.test->service": config.Service{
			Type: "queue",
			Resolve: config.Resolver{
				Type:     "http-post",
				Endpoint: "http://localhost:8080",
				Payload:  []byte("test"),
			},
		},
	}

	queue := NewQueueAdapter(c, resolver)
	id, err := queue.QueueWithOpts(&types.Service{
		Name: "latest.test->service",
		Addr: "http://localhost:8080",
	}, types.Request{}, types.Options{})

	require.NoError(t, err)
	require.Equal(t, id, "test")
}
