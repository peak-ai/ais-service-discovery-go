package backends

import "github.com/peak-ai/ais-service-discovery-go/pkg/types"

// SetQueue sets the given queue adapter to be used
func SetQueue(queueAdapter types.QueueAdapter) Option {
	return func(args *Options) {
		args.QueueAdapter = queueAdapter
	}
}

// SetFunction sets the given function adapter to be used
func SetFunction(functionAdapter types.FunctionAdapter) Option {
	return func(args *Options) {
		args.FunctionAdapter = functionAdapter
	}
}

// SetAutomate sets the given automation adapter to be used
func SetAutomate(automateAdapter types.AutomateAdapter) Option {
	return func(args *Options) {
		args.AutomateAdapter = automateAdapter
	}
}

// SetPubsub sets the given pubsub adapter to be used
func SetPubsub(pubsubAdapter types.PubsubAdapter) Option {
	return func(args *Options) {
		args.PubsubAdapter = pubsubAdapter
	}
}

// SetLocator sets the service discovery adapter to be used
func SetLocator(locator types.Locator) Option {
	return func(args *Options) {
		args.Locator = locator
	}
}

// SetLogger sets the logger service
func SetLogger(logger types.LogAdapter) Option {
	return func(args *Options) {
		args.LogAdapter = logger
	}
}

// SetTracer sets the tracer to be used
func SetTracer(tracer types.TraceAdapter) Option {
	return func(args *Options) {
		args.TraceAdapter = tracer
	}
}

type Option func(*Options)

// Options allows the client to configure
// the integrations.
type Options struct {
	types.QueueAdapter
	types.FunctionAdapter
	types.AutomateAdapter
	types.PubsubAdapter
	types.Locator
	types.LogAdapter
	types.TraceAdapter
}
