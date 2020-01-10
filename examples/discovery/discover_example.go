package main

import (
	"log"

	d "github.com/peak-ai/ais-service-discovery-go"
)

func main() {
	discovery := d.NewDiscovery(d.WithAWSBackend())
	res, err := discovery.Discover("acme-latest.generator->queue")
	if err != nil {
		log.Panic(err)
	}
	log.Println(res)
}
