package discovery

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/peak-ai/ais-service-discovery-go/pkg/automate"
	"github.com/peak-ai/ais-service-discovery-go/pkg/function"
	"github.com/peak-ai/ais-service-discovery-go/pkg/locator"
	"github.com/peak-ai/ais-service-discovery-go/pkg/logger"
	"github.com/peak-ai/ais-service-discovery-go/pkg/parser"
	"github.com/peak-ai/ais-service-discovery-go/pkg/pubsub"
	"github.com/peak-ai/ais-service-discovery-go/pkg/queue"
	"github.com/peak-ai/ais-service-discovery-go/pkg/tracer"
	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
)

// Options allows the client to configure
// the integrations.
type Options struct {
	QueueAdapter
	FunctionAdapter
	AutomateAdapter
	PubsubAdapter
	Locator
	LogAdapter
	TraceAdapter
}

// Option is a function that modifies the Options object
type Option func(*Options)

// SetQueue sets the given queue adapter to be used
func SetQueue(queueAdapter QueueAdapter) Option {
	return func(args *Options) {
		args.QueueAdapter = queueAdapter
	}
}

// SetFunction sets the given function adapter to be used
func SetFunction(functionAdapter FunctionAdapter) Option {
	return func(args *Options) {
		args.FunctionAdapter = functionAdapter
	}
}

// SetAutomate sets the given automation adapter to be used
func SetAutomate(automateAdapter AutomateAdapter) Option {
	return func(args *Options) {
		args.AutomateAdapter = automateAdapter
	}
}

// SetPubsub sets the given pubsub adapter to be used
func SetPubsub(pubsubAdapter PubsubAdapter) Option {
	return func(args *Options) {
		args.PubsubAdapter = pubsubAdapter
	}
}

// SetLocator sets the service discovery adapter to be used
func SetLocator(locator Locator) Option {
	return func(args *Options) {
		args.Locator = locator
	}
}

// SetLogger sets the logger service
func SetLogger(logger LogAdapter) Option {
	return func(args *Options) {
		args.LogAdapter = logger
	}
}

// SetTracer sets the tracer to be used
func SetTracer(tracer TraceAdapter) Option {
	return func(args *Options) {
		args.TraceAdapter = tracer
	}
}

// NewDiscovery creates a new Discover object to communicate with the various services
func NewDiscovery(opts ...Option) *Discover {
	sess := session.Must(session.NewSession())
	args := &Options{
		QueueAdapter:    queue.NewSQSAdapter(sqs.New(sess)),
		FunctionAdapter: function.NewLambdaAdapter(lambda.New(sess)),
		AutomateAdapter: automate.NewSSMAdapter(ssm.New(sess)),
		PubsubAdapter:   pubsub.NewSNSAdapter(sns.New(sess)),
		Locator:         locator.NewCloudmapLocator(servicediscovery.New(sess)),
		LogAdapter:      logger.NewSTDOutAdapter(),
		TraceAdapter:    tracer.NewXrayAdapter(),
	}

	for _, opt := range opts {
		opt(args)
	}

	return &Discover{
		QueueAdapter:    args.QueueAdapter,
		FunctionAdapter: args.FunctionAdapter,
		AutomateAdapter: args.AutomateAdapter,
		PubsubAdapter:   args.PubsubAdapter,
		Locator:         args.Locator,
		LogAdapter:      args.LogAdapter,
		TraceAdapter:    args.TraceAdapter,
	}
}

// QueueAdapter is an interface defining a Queue service
type QueueAdapter interface {

	// Queue a message, return a token or message id
	Queue(service *types.Service, request types.Request, opts types.Options) (string, error)
	Listen(service *types.Service, opts types.Options) (<-chan *types.Response, error)
}

// FunctionAdapter is an interface defining a Serverless Functions service
type FunctionAdapter interface {
	Call(service *types.Service, request types.Request, opts types.Options) (*types.Response, error)
}

// AutomateAdapter is an interface defining a System Management service
type AutomateAdapter interface {
	Execute(service *types.Service, request types.Request, opts types.Options) (*types.Response, error)
}

// PubsubAdapter is an interface defining a PubSub Messaging service
type PubsubAdapter interface {
	Publish(service *types.Service, request types.Request, opts types.Options) error
	Subscribe(service *types.Service, opts types.Options) (<-chan *types.Response, error)
}

// Locator is an interface defining a Service Discovery service
type Locator interface {
	Discover(signature *types.Signature) (*types.Service, error)
}

// LogAdapter is an interface defining a Logging service
type LogAdapter interface {
	Log(service *types.Service, message string)
}

// TraceAdapter is an interface defining a Tracing service
type TraceAdapter interface {
	Trace(service *types.Service)
}

// Discover instance
type Discover struct {
	QueueAdapter
	FunctionAdapter
	AutomateAdapter
	PubsubAdapter
	Locator
	LogAdapter
	TraceAdapter
}

func (d *Discover) discover(signature string) (*types.Service, error) {
	sig, err := parser.ParseAddr(signature)
	if err != nil {
		return nil, err
	}
	return d.Discover(sig)
}

// Request makes synchronous call through the FunctionAdapter
func (d *Discover) Request(signature string, request types.Request, opts types.Options) (*types.Response, error) {
	service, err := d.discover(signature)
	if err != nil {
		return nil, err
	}
	defer d.log(service, fmt.Sprintf("making a request to: %s", signature))
	defer d.trace(service)
	return d.FunctionAdapter.Call(service, request, opts)
}

// Automate calls an infrastructure script through the AutomateAdapter
func (d *Discover) Automate(signature string, request types.Request, opts types.Options) (*types.Response, error) {
	service, err := d.discover(signature)
	if err != nil {
		return nil, err
	}
	defer d.log(service, fmt.Sprintf("running automation: %s", signature))
	defer d.trace(service)
	return d.AutomateAdapter.Execute(service, request, opts)
}

// Publish publishes an asynchronous event through the PubsubAdapter
func (d *Discover) Publish(signature string, request types.Request, opts types.Options) error {
	service, err := d.discover(signature)
	if err != nil {
		return err
	}
	defer d.log(service, fmt.Sprintf("publishing event to: %s", signature))
	defer d.trace(service)
	return d.PubsubAdapter.Publish(service, request, opts)
}

// Subscribe subscribes to an event through the PubsubAdapter
// Warning, not possible with SNS
func (d *Discover) Subscribe(signature string, opts types.Options) (<-chan *types.Response, error) {
	service, err := d.discover(signature)
	if err != nil {
		return nil, err
	}
	defer d.log(service, fmt.Sprintf("subscribed to: %s", signature))
	defer d.trace(service)
	return d.PubsubAdapter.Subscribe(service, opts)
}

// Queue queues a request through the QueueAdapter
func (d *Discover) Queue(signature string, request types.Request, opts types.Options) (string, error) {
	service, err := d.discover(signature)
	if err != nil {
		return "", err
	}
	defer d.log(service, fmt.Sprintf("queued message to: %s", signature))
	defer d.trace(service)
	return d.QueueAdapter.Queue(service, request, opts)
}

// Listen creates a listener channel through the QueueAdapter
func (d *Discover) Listen(signature string, opts types.Options) (<-chan *types.Response, error) {
	service, err := d.discover(signature)
	if err != nil {
		return nil, err
	}
	defer d.log(service, fmt.Sprintf("listening to: %s", signature))
	defer d.trace(service)
	return d.QueueAdapter.Listen(service, opts)
}

// Logs the call that's made, using the given
// log adapter.
func (d *Discover) log(service *types.Service, message string) {
	d.LogAdapter.Log(service, message)
}

// Traces the call that's made, using the given
// trace adapter.
func (d *Discover) trace(service *types.Service) {
	d.TraceAdapter.Trace(service)
}

// Call sends a request to the proper adapter depending on the service type.
// (potentially not needed, as the behavioural methods say)
func (d *Discover) Call(service types.ServiceRequest, opts types.Options) (*types.Response, error) {
	switch service.Service.Type {
	case "function", "lambda":
		return d.FunctionAdapter.Call(service.Service, service.Request, opts)
	case "event", "pubsub":
		return &types.Response{}, d.PubsubAdapter.Publish(service.Service, service.Request, opts)
	case "queue", "sqs":
		token, err := d.QueueAdapter.Queue(service.Service, service.Request, opts)
		return &types.Response{
			Body: []byte(token),
		}, err
	case "script", "ssm", "automation":
		return d.AutomateAdapter.Execute(service.Service, service.Request, opts)
	default:
		// @todo - potentially dangerous default option?
		return d.FunctionAdapter.Call(service.Service, service.Request, opts)
	}
}
