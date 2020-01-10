package locator

import (
	"errors"
	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicediscovery"

	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
)

// CloudmapLocator is an implementation of the Locator using AWS CloudMap
type CloudmapLocator struct {
	client *servicediscovery.ServiceDiscovery
}

// NewCloudmapLocator creates a new instance of CloudmapLocator
func NewCloudmapLocator(client *servicediscovery.ServiceDiscovery) *CloudmapLocator {
	return &CloudmapLocator{client}
}

// Discover searches for a service
func (l *CloudmapLocator) Discover(service *types.Signature) (*types.Service, error) {
	input := &servicediscovery.DiscoverInstancesInput{
		NamespaceName: aws.String(service.Namespace),
		ServiceName:   aws.String(service.Service),
	}

	instanceOutput, err := l.client.DiscoverInstances(input)
	if err != nil {
		return nil, err
	}

	var instance *servicediscovery.HttpInstanceSummary
	instances := instanceOutput.Instances
	for _, i := range instances {
		if i.InstanceId == aws.String(service.Instance) {
			instance = i
		}
	}


	// @todo - 'arn' is AWS specific, consider a more
	// generalised term for this concept
	// also, sometimes people use 'arn' other cases url 'url'.
	location := *instance.Attributes["arn"]
	if location == "" {
		location = *instance.Attributes["url"]
	}

	if location == "" {
		return nil, errors.New("cannot find a url or arn associated with this service")
	}

	t := *instance.Attributes["type"]

	return &types.Service{
		Name: service.Service,
		Addr: location,
		Type: t,
	}, nil
}
