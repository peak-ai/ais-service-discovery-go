package backend

import (
	"log"

	discovery "github.com/peak-ai/ais-service-discovery-go"
	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
	"github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/config"
	"github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/mocked"
	"github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/redirected"
	"github.com/peak-ai/ais-service-discovery-go/plugins/incubator/local-plane/resolvers"
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

// WithMockedBackend -
func WithMockedBackend(config config.Config) discovery.Option {
	logger, _ := zap.NewProduction()
	structuredLogger := &Logger{logger}
	tracer := &FakeTracer{}

	return func(args *discovery.Options) {
		args.TraceAdapter = tracer
		args.LogAdapter = structuredLogger
		args.Locator = mocked.NewLocator(config)
		args.QueueAdapter = mocked.NewQueueAdapter(config)
	}
}

// WithWithRedirectedBackend -
func WithRedirectedBackend(config config.Config) discovery.Option {
	resolver := resolvers.NewResolver()
	return func(args *discovery.Options) {
		args.Locator = mocked.NewLocator(config)
		args.QueueAdapter = redirected.NewQueueAdapter(config, resolver)
	}
}
