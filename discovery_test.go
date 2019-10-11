package discovery

import (
	"log"
	"testing"

	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
	"github.com/stretchr/testify/assert"
)

type mockLocator struct{}

func (l *mockLocator) Discover(service *types.Signature) (*types.Service, error) {
	return &types.Service{}, nil
}

type mockFunctionAdapter struct{}

func (fa *mockFunctionAdapter) Call(
	service *types.Service,
	request types.Request,
	opts types.Options,
) (*types.Response, error) {
	return nil, nil
}

type mockQueueAdapter struct {
	messages chan *types.Response
}

func (qa *mockQueueAdapter) Queue(
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

func (qa *mockQueueAdapter) Listen(
	service *types.Service,
	options types.Options,
) (<-chan *types.Response, error) {
	return qa.messages, nil
}

type mockAutomateAdapter struct{}

func (aa *mockAutomateAdapter) Execute(
	service *types.Service,
	request types.Request,
	opts types.Options,
) (*types.Response, error) {
	return &types.Response{}, nil
}

type mockPubsubAdapter struct{}

func (pa *mockPubsubAdapter) Publish(
	service *types.Service,
	request types.Request,
	opts types.Options,
) error {
	return nil
}

func (pa *mockPubsubAdapter) Subscribe(
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
	)
	discover.Request("service->handler", types.Request{
		Body: []byte(""),
	}, nil)
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
	}, nil)
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
	}, nil)
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
	}, nil)
	assert.NoError(t, err)
	assert.Equal(t, token, "abc123")

	messages, err := discover.Listen("my-service->my-queue", nil)
	log.Println(messages)
	assert.NoError(t, err)
	assert.Equal(t, <-messages, response)
}
