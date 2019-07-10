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

func TestCanCallRequest(t *testing.T) {
	discover := &Discover{
		Locator:         &mockLocator{},
		FunctionAdapter: &mockFunctionAdapter{},
	}
	discover.Request("service->handler", types.Request{
		Body: []byte(""),
	}, nil)
}

func TestCanCallAutomate(t *testing.T) {
	discover := NewDiscovery(
		SetLocator(&mockLocator{}),
		SetAutomate(&mockAutomateAdapter{}),
	)
	_, err := discover.Automate("service", types.Request{
		Body: []byte(""),
	}, nil)
	assert.NoError(t, err)
}

func TestCanCallPubSub(t *testing.T) {
	discover := NewDiscovery(
		SetLocator(&mockLocator{}),
		SetPubsub(&mockPubsubAdapter{}),
	)
	err := discover.Publish("my-topic", types.Request{
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
	)
	token, err := discover.Queue("my-queue", types.Request{
		Body: val,
	}, nil)
	assert.NoError(t, err)
	assert.Equal(t, token, "abc123")

	messages, err := discover.Listen("my-queue", nil)
	log.Println(messages)
	assert.NoError(t, err)
	assert.Equal(t, <-messages, response)
}
