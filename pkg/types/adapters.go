package types

// QueueAdapter is an interface defining a Queue service
type QueueAdapter interface {

	// Queue a message, return a token or message id
	QueueWithOpts(service *Service, request Request, opts Options) (string, error)
	ListenWithOpts(service *Service, opts Options) (<-chan *Response, error)
}

// FunctionAdapter is an interface defining a Serverless Functions service
type FunctionAdapter interface {
	CallWithOpts(service *Service, request Request, opts Options) (*Response, error)
}

// AutomateAdapter is an interface defining a System Management service
type AutomateAdapter interface {
	ExecuteWithOpts(service *Service, request Request, opts Options) (*Response, error)
}

// PubsubAdapter is an interface defining a PubSub Messaging service
type PubsubAdapter interface {
	PublishWithOpts(service *Service, request Request, opts Options) error
	SubscribeWithOpts(service *Service, opts Options) (<-chan *Response, error)
}

// Locator is an interface defining a Service Discovery service
type Locator interface {
	Discover(signature *Signature) (*Service, error)
}

// LogAdapter is an interface defining a Logging service
type LogAdapter interface {
	Log(service *Service, message string)
}

// TraceAdapter is an interface defining a Tracing service
type TraceAdapter interface {
	Trace(service *Service)
}
