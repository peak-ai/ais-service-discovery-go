package pubsub

import (
	"errors"

	"github.com/peak-ai/ais-service-discovery-go/pkg/types"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

// SNSAdapter -
type SNSAdapter struct {
	client *sns.SNS
}

// NewSNSAdapter instance
func NewSNSAdapter(client *sns.SNS) *SNSAdapter {
	return &SNSAdapter{client}
}

// Publish -
func (sa *SNSAdapter) Publish(service *types.Service, request types.Request, opts types.Options) error {
	input := &sns.PublishInput{
		Message:  aws.String(string(request.Body)),
		TopicArn: aws.String(service.Addr),
	}
	_, err := sa.client.Publish(input)
	return err
}

// Subscribe - subscriptions are at a higher, none code level for AWS,
// so we can't subscribe through code as such.
func (sa *SNSAdapter) Subscribe(service *types.Service, opts types.Options) (<-chan *types.Response, error) {
	return nil, errors.New("not valid for SNS")
}
