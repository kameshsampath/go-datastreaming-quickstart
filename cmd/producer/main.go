package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	seeds := []string{os.Getenv("RPK_BROKERS")}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup("my-group"),
		kgo.ConsumeTopics("greetings"),
	)

	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	//Produce
	bctx := context.Background()
	tctx, cancel := context.WithTimeout(bctx, 3*time.Second)
	defer cancel()

	//Async Producer
	var wg sync.WaitGroup
	wg.Add(1)
	kr := &kgo.Record{
		Key:   []byte("1"),
		Value: []byte("Hello World."),
		Topic: "greetings",
	}

	client.Produce(tctx, kr, func(r *kgo.Record, err error) {
		if err != nil {
			log.Fatal(err)
		}
		defer wg.Done()
		fmt.Printf("Saved record to Partition %d, Offset %d", kr.Partition, kr.Offset)
	})

	wg.Wait()
}
