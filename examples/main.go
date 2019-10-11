package main

import (
	"log"

	d "github.com/peak-ai/ais-service-discovery-go"
	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
)

func main() {
	discovery := d.NewDiscovery()
	res, err := discovery.Request("acme-prod.scheduler->create-job", types.Request{
		Body: []byte(`{ "frequency": "* * * * *", "type": "test" }`),
	}, nil)
	if err != nil {
		log.Panic(err)
	}
	log.Println(string(res.Body))

	messages, err := discovery.Listen("acme-prod.service-discovery->test-queue", nil)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		discovery.Queue("acme-prod.service-discovery->test-queue", types.Request{
			Body: []byte(`{ "message": "hello" }`),
		}, nil)
	}()

	go func() {
		discovery.Queue("acme-prod.service-discovery->test-queue", types.Request{
			Body: []byte(`{ "message": "again" }`),
		}, nil)
	}()

	for message := range messages {
		log.Println(string(message.Body))
	}
}
