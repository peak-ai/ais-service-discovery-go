package main

import "github.com/peak-ai/ais-service-discovery-go/pkg/backends/mocked"

func main() {
	config := mocked.Config{}
	sd := mocked.Factory(config)
}
