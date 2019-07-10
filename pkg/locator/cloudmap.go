package locator

import (
	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicediscovery"

	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
)

// CloudmapLocator -
type CloudmapLocator struct {
	client *servicediscovery.ServiceDiscovery
}

// NewCloudmapLocator -
func NewCloudmapLocator(client *servicediscovery.ServiceDiscovery) *CloudmapLocator {
	return &CloudmapLocator{client}
}

// Discover a service
func (l *CloudmapLocator) Discover(service *types.Signature) (*types.Service, error) {
	input := &servicediscovery.DiscoverInstancesInput{
		NamespaceName: aws.String(service.Namespace),
		ServiceName:   aws.String(service.Service),
	}

	if service.Handler != "" {
		input.QueryParameters = make(map[string]*string)
		input.QueryParameters["handler"] = aws.String(service.Handler)
	}

	instanceOutput, err := l.client.DiscoverInstances(input)
	if err != nil {
		return nil, err
	}

	// @todo - we might want to actually deal with
	// multiple instances at some point.
	instance := instanceOutput.Instances[0]

	// @todo - 'arn' is AWS specific, consider a more
	// generalisd term for this concept
	location := *instance.Attributes["arn"]
	t := *instance.Attributes["type"]

	return &types.Service{
		Name: service.Service,
		Addr: location,
		Type: t,
	}, nil
}
