package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/twmb/franz-go/pkg/kgo"
)

// Produce produces the data to a topic
func Produce(c echo.Context) error {
	var request Request
	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	seeds := c.Get(SEED_CONTEXT_KEY).(Seeds)
	client, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
	)

	if err != nil {
		sugar.Fatal(err)
	}
	defer client.Close()

	//Produce
	ctx := context.Background()
	tctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	//Sync Producer
	//TODO do async produce ??
	kr := &kgo.Record{
		Key:   []byte(request.Key),
		Value: []byte(request.Value),
		Topic: request.Topic,
	}

	if p := request.Partition; p != "" {
		p, err := strconv.Atoi(p)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad partition")
		}
		kr.Partition = int32(p)
	}

	if err := client.ProduceSync(tctx, kr).FirstErr(); err != nil {
		c.Logger().Errorf("error sending message to topic %s,%s", request.Topic, err)
		return err
	}

	sugar.Infow("Successfully posted message", "topic", request.Topic, "Key", request.Key, "Value", request.Value)

	return nil

}
