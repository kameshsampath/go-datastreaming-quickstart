package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Result struct {
	Record *kgo.Record
	Errors []kgo.FetchError
}

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

	ch := make(chan Result)
	go func() {
		pollAndPrint(client, ch)
	}()

	for {
		select {
		case r := <-ch:
			{
				if errs := r.Errors; len(errs) > 0 {
					log.Fatal(errs)
				} else {
					fmt.Printf("Key:%s, Value:%s\n", r.Record.Key, r.Record.Value)
				}
			}
		}
	}

}

func pollAndPrint(client *kgo.Client, ch chan Result) {
	log.Println("Started to poll topic until error")
	bctx := context.Background()

	//Consumer
	for {
		fetches := client.PollFetches(bctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			ch <- Result{
				Errors: errs,
			}
		}

		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			for _, r := range p.Records {
				ch <- Result{
					Record: r,
				}
			}
		})

	}
}
