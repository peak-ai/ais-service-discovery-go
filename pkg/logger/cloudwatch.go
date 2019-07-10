package logger

import (
	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
)

// NewCloudwatchLogger - 
func NewCloudwatchLogger() *CloudwatchLogger {
	return &CloudwatchLogger{}
}

// CloudwatchLogger - 
type CloudwatchLogger struct {}

// Log - 
func (cw *CloudwatchLogger) Log(service *types.Service, message string) {
	return
}