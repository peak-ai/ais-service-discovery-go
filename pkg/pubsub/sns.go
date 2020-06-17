package pubsub

import (
	"errors"

	"github.com/peak-ai/ais-service-discovery-go/pkg/types"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

// SNSAdapter is an implementation of a PubsubAdapter using AWS SNS
type SNSAdapter struct {
	client *sns.SNS
}

// NewSNSAdapter creates a new SNSAdapter instance
func NewSNSAdapter(client *sns.SNS) *SNSAdapter {
	return &SNSAdapter{client}
}

func (sa *SNSAdapter) parseOpts(opts types.Options) map[string]*sns.MessageAttributeValue {
	atts := make(map[string]*sns.MessageAttributeValue, 0)
	for key, val := range opts {
		attributeValue, ok := val.(*sns.MessageAttributeValue)
		if ok {
			atts[key] = attributeValue
		}
	}

	return atts
}

// Publish publishes an event to a queue
func (sa *SNSAdapter) Publish(service *types.Service, request types.Request) error {
	return sa.PublishWithOpts(service, request, types.Options{})
}

// PublishWithOpts takes the generic options type, converts to 'MessageAttributes'
func (sa *SNSAdapter) PublishWithOpts(service *types.Service, request types.Request, opts types.Options) error {
	input := &sns.PublishInput{
		Message:  aws.String(string(request.Body)),
		TopicArn: aws.String(service.Addr),
	}
	
	if len(opts) > 0 {
		atts := sa.parseOpts(opts)
		input.SetMessageAttributes(atts)
	}

	_, err := sa.client.Publish(input)
	return err
}

// Subscribe is not implemented
// (subscriptions are at a higher, none code level for AWS,
// so we can't subscribe through code as such)
func (sa *SNSAdapter) Subscribe(service *types.Service) (<-chan *types.Response, error) {
	return sa.SubscribeWithOpts(service, types.Options{})
}

// SubscribeWithOpts is not implemented
// (subscriptions are at a higher, none code level for AWS,
// so we can't subscribe through code as such)
func (sa *SNSAdapter) SubscribeWithOpts(service *types.Service, opts types.Options) (<-chan *types.Response, error) {
	return nil, errors.New("not valid for SNS")
}
