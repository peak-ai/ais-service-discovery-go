# Cloud Application Framework

This library aims to expose a common set of interfaces for handling common cloud platform tasks. Such as queuing messages, publishing events, calling cloud functions etc.

This library is AWS centric, however, can by modified and extended to support others. Using the interfaces and configuration options shown below.

Being AWS centric, the default options are:

- Locator / Service Discovery: AWS Cloudmap
- Request: AWS Lambda
- Pubsub: AWS SNS
- Queue: AWS SQS
- Automate: AWS SSM

## Use with default settings

```golang
func main() {
  d := discover.NewDiscovery()

  token, err := d.Queue("ais-latest.my-queue", types.Request{
    Body: []byte("{}"),
  }, nil)
  ...
}
```

## Use with custom integrations

```golang
func main() {
  ...
  d := discover.NewDiscovery(
    discover.SetQueue(NewKafkaAdapter(kafkaClient)),
    discover.SetPubsub(NewNATSAdapter(natsClient)),
    discover.SetLocator(NewConsulLocator(consul)),
  )
}
```

### Request

```golang
d := discover.NewDiscovery()
d.Request("my-namespace.my-service->my-function", types.Request{
  Body: []byte("{ \"hello\": \"world\" }"),
}, opts)
```

### Queue

```golang
d := discover.NewDiscovery()
d.Queue("my-namespace.my-queue", types.Request{
  Body: jsonString,
}, nil)

go func() {
  messages, err := d.Listen("my-namespace.my-queue", opts)
  for message := range message {
    log.Println(string(message.Body))
  }
}()
```

### Pubsub

```golang
d := discovery.NewDiscovery()
d.Publish("my-namespace.my-event", types.Request{
  Body: jsonEvent,
})
```

### Automate

```golang
d := discovery.NewDiscovery()
d.Automate("my-namespace.my-script", types.Request{
  Body: jsonEvent,
})
```

## Custom integrations

You can customise the behaviour and create your own integrations by conforming to the following interfaces, and use the `SetLocator`, `SetQueue`, `SetPubsub`, `SetAutomate` and `SetFunction` methods when creating an instance of the Discovery library.

### Locator interface

```golang
Discover(signature *types.Signature) (*types.Service, error)
```

### Queue interface

```golang
// QueueAdapter -
type QueueAdapter interface {

	// Queue a message, return a token or message id
	Queue(service *types.Service, request types.Request, opts types.Options) (string, error)
	Listen(service *types.Service, opts types.Options) (<-chan *types.Response, error)
}
```

### Function interface

```golang
// FunctionAdapter -
type FunctionAdapter interface {
	Call(service *types.Service, request types.Request, opts types.Options) (*types.Response, error)
}
```

### AutomateAdapter interface

```golang
// AutomateAdapter -
type AutomateAdapter interface {
	Execute(service *types.Service, request types.Request, opts types.Options) (*types.Response, error)
}
```

### PubsubAdapter interface

```golang
// PubsubAdapter -
type PubsubAdapter interface {
	Publish(service *types.Service, request types.Request, opts types.Options) error
	Subscribe(service *types.Service, opts types.Options) (<-chan *types.Response, error)
}
```

_Acknowledgements:_ inspired by the amazing work at [Micro](https://github.com/micro/micro)
