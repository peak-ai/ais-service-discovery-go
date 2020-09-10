package mocked

import (
	"log"

	discovery "github.com/peak-ai/ais-service-discovery-go"
	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
	"go.uber.org/zap"
)

// FakeTracer -
type FakeTracer struct{}

// Trace -
func (t *FakeTracer) Trace(service *types.Service) {
	log.Println(service)
}

// Logger -
type Logger struct {
	logger *zap.Logger
}

// Log -
func (l *Logger) Log(service *types.Service, message string) {
	logger := l.logger.With(
		zap.String("addr", service.Addr),
		zap.String("name", service.Name),
		zap.String("type", service.Type),
	)
	logger.Info(message)
}

// Resolver is the action to perform instead of the
// original action. I.e http post request
type Resolver struct {
	Type         string `json:"type"`
	Endpoint     string `json:"endpoint"`
	Payload      []byte `json:"payload"`
	MockResponse []byte `json:"mock_response"`
}

// Service type is the original type,
// i.e. queue, event etc.
// Resolver is the action to perform instead.
type Service struct {
	Type    string   `json:"type"`
	Resolve Resolver `json:"resolve"`
}

// Config the key is a service address, and the value
// is the metadata used to figure out what to do when
// that service is called.
type Config map[string]Service

// WithMockedBackend -
func Factory(config Config) discovery.Option {
	logger, _ := zap.NewProduction()
	structuredLogger := &Logger{logger}
	tracer := &FakeTracer{}

	return func(args *discovery.Options) {
		args.TraceAdapter = tracer
		args.LogAdapter = structuredLogger
		args.Locator = NewLocator(config)
		args.QueueAdapter = NewQueueAdapter(config)
	}
}
