package automate

import (
	"encoding/json"

	"github.com/peak-ai/ais-service-discovery-go/pkg/types"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// SSMAdapter is a AutomateAdapter for AWS, using the Systems Manager service
type SSMAdapter struct {
	client *ssm.SSM
}

// NewSSMAdapter creates a new instance of the SSMAdapter
func NewSSMAdapter(client *ssm.SSM) *SSMAdapter {
	return &SSMAdapter{client}
}

// Execute executes an SSM document, with arguments and returns the execution ID
// as the response body.
func (sa *SSMAdapter) Execute(service *types.Service, request types.Request) (*types.Response, error) {
	return sa.ExecuteWithOpts(service, request, types.Options{})
}

// ExecuteWithOpts - executes an SSM document, with arguments and options.  Returns the execution ID
// as the response body.
func (sa *SSMAdapter) ExecuteWithOpts(service *types.Service, request types.Request, opts types.Options) (*types.Response, error) {
	var args map[string][]*string
	if err := json.Unmarshal(request.Body, &args); err != nil {
		return nil, err
	}
	input := &ssm.StartAutomationExecutionInput{
		DocumentName: aws.String(string(service.Addr)),
		Parameters:   args,
	}
	output, err := sa.client.StartAutomationExecution(input)
	return &types.Response{
		Body: []byte(*output.AutomationExecutionId),
	}, err
}
