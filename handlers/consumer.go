package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/twmb/franz-go/pkg/kgo"
)

// Consume methods consumes the message from the group
// and stream it
func Consume(c echo.Context) error {
	var query Query
	if err := c.Bind(&query); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)

	c.Logger().Infof("Streaming messages for topic, %s", query.Topic)

	seeds := c.Get(SEED_CONTEXT_KEY).(Seeds)

	groupID := fmt.Sprintf("%s-echo-group", query.Topic)

	if query.Group == "new" {
		groupID = fmt.Sprintf("%s-echo-group-%s", query.Topic, uuid.New())
	}

	client, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumeTopics(query.Topic),
		kgo.ConsumerGroup(groupID),
	)

	if err != nil {
		c.Logger().Errorf("error building streaming client %s", err)
		return err
	}
	defer client.Close()

	enc := json.NewEncoder(c.Response())

	ch := make(chan Response)
	go func() {
		poll(client, ch)
	}()

	for {
		select {
		case r := <-ch:
			{
				if err := enc.Encode(r); err != nil {
					return err
				}
				c.Response().Flush()
			}
		}
	}
}

// polls topic for messages and sends the record to the channel
// if err then error is send as part of channel
func poll(client *kgo.Client, ch chan Response) {
	ctx := context.Background()

	for {
		fetches := client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			ch <- Response{
				Errors: errs,
			}
		}

		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			for _, r := range p.Records {
				ch <- Response{
					Key:       string(r.Key),
					Value:     string(r.Value),
					Partition: fmt.Sprintf("%d", r.Partition),
					Offset:    fmt.Sprintf("%d", r.Offset),
					Timestamp: r.Timestamp.Format(time.RFC850),
				}
			}
		})

	}
}
