package main

import (
	"github.com/peak-ai/ais-service-discovery-go/pkg/parser"
	"log"

	d "github.com/peak-ai/ais-service-discovery-go"
)

func main() {
	discovery := d.NewDiscovery(d.WithAWSBackend())
	signature, err := parser.ParseAddr("acme-prod.my-service->queue")
	res, err := discovery.Discover(signature)
	if err != nil {
		log.Panic(err)
	}
	log.Println(res)
}
