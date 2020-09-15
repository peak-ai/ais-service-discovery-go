package discovery

import (
	"log"
	"testing"

	"github.com/peak-ai/ais-service-discovery-go/types"
	"github.com/stretchr/testify/assert"
)

type mockLocator struct{}

func (l *mockLocator) Discover(service *types.Signature) (*types.Service, error) {
	return &types.Service{}, nil
}

type mockFunctionAdapter struct{}

func (fa *mockFunctionAdapter) CallWithOpts(
	service *types.Service,
	request types.Request,
	opts types.Options,
) (*types.Response, error) {
	return nil, nil
}

type mockQueueAdapter struct {
	messages chan *types.Response
}

func (qa *mockQueueAdapter) QueueWithOpts(
	service *types.Service,
	request types.Request,
	opts types.Options,
) (string, error) {
	response := &types.Response{
		Body: request.Body,
	}
	go func() { qa.messages <- response }()
	return "abc123", nil
}

func (qa *mockQueueAdapter) ListenWithOpts(
	service *types.Service,
	options types.Options,
) (<-chan *types.Response, error) {
	return qa.messages, nil
}

type mockAutomateAdapter struct{}

func (aa *mockAutomateAdapter) ExecuteWithOpts(
	service *types.Service,
	request types.Request,
	opts types.Options,
) (*types.Response, error) {
	return &types.Response{}, nil
}

type mockPubsubAdapter struct{}

func (pa *mockPubsubAdapter) PublishWithOpts(
	service *types.Service,
	request types.Request,
	opts types.Options,
) error {
	return nil
}

func (pa *mockPubsubAdapter) SubscribeWithOpts(
	service *types.Service,
	opts types.Options,
) (<-chan *types.Response, error) {
	channel := make(chan *types.Response)
	return channel, nil
}

type mockLogger struct{}

func (ml *mockLogger) Log(service *types.Service, message string) {}

type mockTracer struct{}

func (mt *mockTracer) Trace(service *types.Service) {}

func TestCanCallRequest(t *testing.T) {
	discover := NewDiscovery(
		SetLocator(&mockLocator{}),
		SetAutomate(&mockAutomateAdapter{}),
		SetLogger(&mockLogger{}),
		SetTracer(&mockTracer{}),
		SetFunction(&mockFunctionAdapter{}),
	)
	discover.Request("service->handler", types.Request{
		Body: []byte(""),
	})
}

func TestCanCallAutomate(t *testing.T) {
	discover := NewDiscovery(
		SetLocator(&mockLocator{}),
		SetAutomate(&mockAutomateAdapter{}),
		SetLogger(&mockLogger{}),
		SetTracer(&mockTracer{}),
	)
	_, err := discover.Automate("service->handler", types.Request{
		Body: []byte(""),
	})
	assert.NoError(t, err)
}

func TestCanCallPubSub(t *testing.T) {
	discover := NewDiscovery(
		SetLocator(&mockLocator{}),
		SetPubsub(&mockPubsubAdapter{}),
		SetLogger(&mockLogger{}),
		SetTracer(&mockTracer{}),
	)
	err := discover.Publish("my-service->my-topic", types.Request{
		Body: []byte(""),
	})
	assert.NoError(t, err)
}

func TestCanCallQueue(t *testing.T) {
	val := []byte("{}")
	response := &types.Response{
		Body: val,
	}
	discover := NewDiscovery(
		SetLocator(&mockLocator{}),
		SetQueue(&mockQueueAdapter{
			messages: make(chan *types.Response),
		}),
		SetLogger(&mockLogger{}),
		SetTracer(&mockTracer{}),
	)
	token, err := discover.Queue("my-service->my-queue", types.Request{
		Body: val,
	})
	assert.NoError(t, err)
	assert.Equal(t, token, "abc123")

	messages, err := discover.Listen("my-service->my-queue")
	log.Println(messages)
	assert.NoError(t, err)
	assert.Equal(t, <-messages, response)
}
