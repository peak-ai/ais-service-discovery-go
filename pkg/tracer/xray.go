package tracer

import (
	"context"

	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
	"github.com/aws/aws-xray-sdk-go/xray"
)

// XrayTracer -
type XrayTracer struct{}

// NewXrayAdapter -
func NewXrayAdapter() *XrayTracer {
	return &XrayTracer{}
}

// Trace -
func (xt *XrayTracer) Trace(service *types.Service) {
	_, seg := xray.BeginSegment(context.Background(), service.Name)
	seg.Close(nil)
}
