package resolvers

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
	"github.com/peak-ai/local-plane/config"
	"github.com/pkg/errors"
)

// NewHttpResolver -
func NewHttpResolver(client *http.Client) *HttpResolver {
	return &HttpResolver{client}
}

// HttpResolver -
type HttpResolver struct {
	client *http.Client
}

// Get -
func (r *HttpResolver) Get(resolver config.Resolver) (*types.Response, error) {
	res, err := r.client.Get(resolver.Endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching endpoint")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error unwrapping response body")
	}

	return &types.Response{
		Body: body,
	}, nil
}

// Post -
func (r *HttpResolver) Post(resolver config.Resolver) (*types.Response, error) {
	requestBody := bytes.NewReader(resolver.Payload)
	req, err := http.NewRequest("POST", resolver.Endpoint, requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "error creating post request")
	}

	res, err := r.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error making post request")
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	return &types.Response{
		Body: body,
	}, nil
}
