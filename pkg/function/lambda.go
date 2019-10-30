package function

import (
	"github.com/peak-ai/ais-service-discovery-go/pkg/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// LambdaAdapter is an implentation of the FunctionAdapter using AWS Lambda
type LambdaAdapter struct {
	client *lambda.Lambda
}

// NewLambdaAdapter creates a new instance of the LambdaAdapter
func NewLambdaAdapter(client *lambda.Lambda) *LambdaAdapter {
	return &LambdaAdapter{client}
}

// Call exeutes a lambda function
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
