package function

import (
	"github.com/peak-ai/ais-service-discovery-go/pkg/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// LambdaAdapter -
type LambdaAdapter struct {
	client *lambda.Lambda
}

// NewLambdaAdapter -
func NewLambdaAdapter(client *lambda.Lambda) *LambdaAdapter {
	return &LambdaAdapter{client}
}

// Call a lambda function
func (la *LambdaAdapter) Call(service *types.Service, request types.Request, opts types.Options) (*types.Response, error) {
	input := &lambda.InvokeInput{
		FunctionName: aws.String(service.Addr),
		Payload:      request.Body,
	}

	result, err := la.client.Invoke(input)
	return &types.Response{
		Body: result.Payload,
	}, err
}
