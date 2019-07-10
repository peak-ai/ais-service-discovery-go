package queue

import (
	"github.com/peak-ai/ais-service-discovery-go/pkg/types"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSAdapter -
type SQSAdapter struct {
	client *sqs.SQS
}

// NewSQSAdapter -
func NewSQSAdapter(client *sqs.SQS) *SQSAdapter {
	return &SQSAdapter{client}
}

// Queue a message
func (qa *SQSAdapter) Queue(service *types.Service, request types.Request, opts types.Options) (string, error) {
	input := &sqs.SendMessageInput{
		MessageBody: aws.String(string(request.Body)),
		QueueUrl:    aws.String(service.Addr),
	}
	output, err := qa.client.SendMessage(input)
	return *output.MessageId, err
}

// Listen for messages
func (qa *SQSAdapter) Listen(service *types.Service, opts types.Options) (<-chan *types.Response, error) {
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
