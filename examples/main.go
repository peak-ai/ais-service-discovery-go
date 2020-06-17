package main

import (
	"log"

	d "github.com/peak-ai/ais-service-discovery-go"
	"github.com/peak-ai/ais-service-discovery-go/pkg/types"
)

func main() {
	discovery := d.NewDiscovery(d.WithAWSBackend())
	res, err := discovery.Request("acme-prod.scheduler->create-job", types.Request{
		Body: []byte(`{ "frequency": "* * * * *", "type": "test" }`),
	})
	if err != nil {
		log.Panic(err)
	}
	log.Println(string(res.Body))

	messages, err := discovery.Listen("acme-prod.service-discovery->test-queue")
	if err != nil {
		log.Panic(err)
	}

	go func() {
		discovery.Queue("acme-prod.service-discovery->test-queue", types.Request{
			Body: []byte(`{ "message": "hello" }`),
		})
	}()

	go func() {
		discovery.Queue("acme-prod.service-discovery->test-queue", types.Request{
			Body: []byte(`{ "message": "again" }`),
		})
	}()

	for message := range messages {
		log.Println(string(message.Body))
	}
}
