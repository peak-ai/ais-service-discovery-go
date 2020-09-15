package logger

import (
	"github.com/peak-ai/ais-service-discovery-go/types"
)

// NewCloudwatchLogger creates a new instance of the CloudwatchLogger
func NewCloudwatchLogger() *CloudwatchLogger {
	return &CloudwatchLogger{}
}

// CloudwatchLogger is an implementation of the LogAdapter using AWS CloudWatch
type CloudwatchLogger struct{}

// Log logs the mesage
func (cw *CloudwatchLogger) Log(service *types.Service, message string) {
	return
}
