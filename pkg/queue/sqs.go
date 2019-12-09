package queue

import (
	"github.com/peak-ai/ais-service-discovery-go/pkg/types"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSAdapter is an implementation of a QueueAdapter using AWS SQS
type SQSAdapter struct {
	client *sqs.SQS
}

// NewSQSAdapter creates a new implementation of a SQSAdapter
func NewSQSAdapter(client *sqs.SQS) *SQSAdapter {
	return &SQSAdapter{client}
}

// Queue queues a message
func (qa *SQSAdapter) Queue(service *types.Service, request types.Request) (string, error) {
	return qa.QueueWithOpts(service, request, types.Options{})
}

// QueueWithOpts queues a message, with options.
func (qa *SQSAdapter) QueueWithOpts(service *types.Service, request types.Request, opts types.Options) (string, error) {
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(string(request.Body)),
		QueueUrl:    aws.String(service.Addr),
	}
	output, err := qa.client.SendMessage(input)
	return *output.MessageId, err
}

// Listen listens for messages
func (qa *SQSAdapter) Listen(service *types.Service) (<-chan *types.Response, error) {
	return qa.ListenWithOpts(service, types.Options{})
}

// ListenWithOpts listens for messages, with options
func (qa *SQSAdapter) ListenWithOpts(service *types.Service, opts types.Options) (<-chan *types.Response, error) {
	rchan := make(chan *types.Response)
	input := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(service.Addr),
	}
	go func() {
		for {
			res, err := qa.client.ReceiveMessage(input)
			if res == nil {
				continue
			}

			if err != nil {
				rchan <- &types.Response{Error: err}
				continue
			}

			for _, msg := range res.Messages {
				rchan <- &types.Response{
					Body: []byte(*msg.Body),
				}

				// @todo - handle error here...
				qa.client.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      aws.String(service.Addr),
					ReceiptHandle: msg.ReceiptHandle,
				})
			}
		}
	}()
	return rchan, nil
}
