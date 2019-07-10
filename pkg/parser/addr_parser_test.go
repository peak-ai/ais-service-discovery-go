package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanParseAFunction(t *testing.T) {
	svc := "my-service->my-handler"
	signature, err := ParseAddr(svc)
	assert.NoError(t, err)
	assert.Equal(t, signature.Service, "my-service")
	assert.Equal(t, signature.Namespace, "default")
	assert.Equal(t, signature.Handler, "my-handler")
}

func TestCanParseSvcWithNamespace(t *testing.T) {
	svc := "my-namespace.my-service->my-handler"
	signature, err := ParseAddr(svc)
	assert.NoError(t, err)
	assert.Equal(t, signature.Service, "my-service")
	assert.Equal(t, signature.Namespace, "my-namespace")
	assert.Equal(t, signature.Handler, "my-handler")
}
