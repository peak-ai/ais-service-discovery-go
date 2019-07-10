package parser

import (
	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
	"strings"
)

// ParseAddr -
// @todo - this is fairly crude, could be improved.
func ParseAddr(signature string) (*types.Signature, error) {
	namespace := ""
	service := ""
	handler := ""
	parts := strings.Split(signature, ".")
	if len(parts) == 2 {
		namespace = parts[0]
		service = parts[1]
	} else {
		namespace = "default"
		service = parts[0]
	}
	if strings.Contains(service, "->") {
		svcParts := strings.Split(service, "->")
		service = svcParts[0]
		handler = svcParts[1]
	}
	return &types.Signature{
		Namespace: namespace,
		Service:   service,
		Handler:   handler,
	}, nil
}
