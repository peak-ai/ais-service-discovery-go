package tracer

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/peak-ai/ais-service-discovery-go/types"
)

// XrayTracer is an implementation of a TraceAdapter using AWS Xray
type XrayTracer struct{}

// NewXrayAdapter creates a new instance of a XrayTracer
func NewXrayAdapter() *XrayTracer {
	return &XrayTracer{}
}

// Trace traces an event
func (xt *XrayTracer) Trace(service *types.Service) {
	_, seg := xray.BeginSegment(context.Background(), service.Name)
	seg.Close(nil)
}
