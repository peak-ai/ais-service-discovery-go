package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanParseAFunction(t *testing.T) {
	svc := "my-service->my-instance"
	signature, err := ParseAddr(svc)
	assert.NoError(t, err)
	assert.Equal(t, signature.Service, "my-service")
	assert.Equal(t, signature.Namespace, "default")
	assert.Equal(t, signature.Instance, "my-instance")
}

func TestCanParseSvcWithNamespace(t *testing.T) {
	svc := "my-namespace.my-service->my-instance"
	signature, err := ParseAddr(svc)
	assert.NoError(t, err)
	assert.Equal(t, signature.Service, "my-service")
	assert.Equal(t, signature.Namespace, "my-namespace")
	assert.Equal(t, signature.Instance, "my-instance")
}
