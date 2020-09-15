package logger

import (
	"os"

	"github.com/peak-ai/ais-service-discovery-go/types"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

// STDOutLogger logs formatted logs to std out.
type STDOutLogger struct{}

// NewSTDOutAdapter creates a new instance of the STDOutAdapter
func NewSTDOutAdapter() *STDOutLogger {
	return &STDOutLogger{}
}

// Log writes to std out
func (so *STDOutLogger) Log(service *types.Service, message string) {
	log.WithFields(log.Fields{
		"service":  service.Name,
		"location": service.Addr,
		"type":     service.Type,
	}).Info(message)
}
