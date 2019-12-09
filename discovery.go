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

// WithAWSBackend initializes the Discovery object with default AWS services.
// Override these services by using the Set methods.
func WithAWSBackend() Option {
	return func(args *Options) {
		sess := session.Must(session.NewSession())
		args.QueueAdapter = queue.NewSQSAdapter(sqs.New(sess))
		args.FunctionAdapter = function.NewLambdaAdapter(lambda.New(sess))
		args.AutomateAdapter = automate.NewSSMAdapter(ssm.New(sess))
		args.PubsubAdapter = pubsub.NewSNSAdapter(sns.New(sess))
		args.Locator = locator.NewCloudmapLocator(servicediscovery.New(sess))
		args.LogAdapter = logger.NewSTDOutAdapter()
		args.TraceAdapter = tracer.NewXrayAdapter()
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
}

// NewDiscovery -
func NewDiscovery(opts ...Option) *Discover {
	args := &Options{}
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
	QueueWithOpts(service *types.Service, request types.Request, opts types.Options) (string, error)
	ListenWithOpts(service *types.Service, opts types.Options) (<-chan *types.Response, error)
}

// FunctionAdapter is an interface defining a Serverless Functions service
type FunctionAdapter interface {
	CallWithOpts(service *types.Service, request types.Request, opts types.Options) (*types.Response, error)
}

// AutomateAdapter is an interface defining a System Management service
type AutomateAdapter interface {
	ExecuteWithOpts(service *types.Service, request types.Request, opts types.Options) (*types.Response, error)
}

// PubsubAdapter is an interface defining a PubSub Messaging service
type PubsubAdapter interface {
	PublishWithOpts(service *types.Service, request types.Request, opts types.Options) error
	SubscribeWithOpts(service *types.Service, opts types.Options) (<-chan *types.Response, error)
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
func (d *Discover) Request(signature string, request types.Request) (*types.Response, error) {
	return d.RequestWithOpts(signature, request, types.Options{})
}

// RequestWithOpts makes synchronous call through the FunctionAdapter, with options
func (d *Discover) RequestWithOpts(signature string, request types.Request, opts types.Options) (*types.Response, error) {
	service, err := d.discover(signature)
	if err != nil {
		return nil, err
	}
	defer d.log(service, fmt.Sprintf("making a request to: %s", signature))
	defer d.trace(service)
	return d.FunctionAdapter.CallWithOpts(service, request, opts)
}

// Automate calls an infrastructure script through the AutomateAdapter
func (d *Discover) Automate(signature string, request types.Request) (*types.Response, error) {
	return d.AutomateWithOpts(signature, request, types.Options{})
}

// AutomateWithOpts calls an infrastructure script through the AutomateAdapter, with options
func (d *Discover) AutomateWithOpts(signature string, request types.Request, opts types.Options) (*types.Response, error) {
	service, err := d.discover(signature)
	if err != nil {
		return nil, err
	}
	defer d.log(service, fmt.Sprintf("running automation: %s", signature))
	defer d.trace(service)
	return d.AutomateAdapter.ExecuteWithOpts(service, request, opts)
}

// Publish publishes an asynchronous event through the PubsubAdapter
func (d *Discover) Publish(signature string, request types.Request) error {
	return d.PublishWithOpts(signature, request, types.Options{})
}

// PublishWithOpts - publishes an asynchronous event through the PubsubAdapter, with options
func (d *Discover) PublishWithOpts(signature string, request types.Request, opts types.Options) error {
	service, err := d.discover(signature)
	if err != nil {
		return err
	}
	defer d.log(service, fmt.Sprintf("publishing event to: %s", signature))
	defer d.trace(service)
	return d.PubsubAdapter.PublishWithOpts(service, request, opts)
}

// Subscribe subscribes to an event through the PubsubAdapter
// Warning, not possible with SNS
func (d *Discover) Subscribe(signature string) (<-chan *types.Response, error) {
	return d.SubscribeWithOpts(signature, types.Options{})
}

// SubscribeWithOpts subscribes to an event through the PubsubAdapter, with options
// Warning, not possible with SNS
func (d *Discover) SubscribeWithOpts(signature string, opts types.Options) (<-chan *types.Response, error) {
	service, err := d.discover(signature)
	if err != nil {
		return nil, err
	}
	defer d.log(service, fmt.Sprintf("subscribed to: %s", signature))
	defer d.trace(service)
	return d.PubsubAdapter.SubscribeWithOpts(service, opts)
}

// Queue queues a request through the QueueAdapter
func (d *Discover) Queue(signature string, request types.Request) (string, error) {
	return d.QueueWithOpts(signature, request, types.Options{})
}

// QueueWithOpts queues a request through the QueueAdapter, with options
func (d *Discover) QueueWithOpts(signature string, request types.Request, opts types.Options) (string, error) {
	service, err := d.discover(signature)
	if err != nil {
		return "", err
	}
	defer d.log(service, fmt.Sprintf("queued message to: %s", signature))
	defer d.trace(service)
	return d.QueueAdapter.QueueWithOpts(service, request, opts)
}

// Listen creates a listener channel through the QueueAdapter
func (d *Discover) Listen(signature string) (<-chan *types.Response, error) {
	return d.ListenWithOpts(signature, types.Options{})
}

// ListenWithOpts creates a listener channel through the QueueAdapter, with options
func (d *Discover) ListenWithOpts(signature string, opts types.Options) (<-chan *types.Response, error) {
	service, err := d.discover(signature)
	if err != nil {
		return nil, err
	}
	defer d.log(service, fmt.Sprintf("listening to: %s", signature))
	defer d.trace(service)
	return d.QueueAdapter.ListenWithOpts(service, opts)
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
	return d.CallWithOpts(service, types.Options{})
}

// CallWithOpts sends a request to the proper adapter depending on the service
// type, with options.
// (potentially not needed, as the behavioural methods say)
func (d *Discover) CallWithOpts(service types.ServiceRequest, opts types.Options) (*types.Response, error) {
	switch service.Service.Type {
	case "function", "lambda":
		return d.FunctionAdapter.CallWithOpts(service.Service, service.Request, opts)
	case "event", "pubsub":
		return &types.Response{}, d.PubsubAdapter.PublishWithOpts(service.Service, service.Request, opts)
	case "queue", "sqs":
		token, err := d.QueueAdapter.QueueWithOpts(service.Service, service.Request, opts)
		return &types.Response{
			Body: []byte(token),
		}, err
	case "script", "ssm", "automation":
		return d.AutomateAdapter.ExecuteWithOpts(service.Service, service.Request, opts)
	default:
		// @todo - potentially dangerous default option?
		return d.FunctionAdapter.CallWithOpts(service.Service, service.Request, opts)
	}
}
