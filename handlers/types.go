package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

// SEED_CONTEXT_KEY is the echo.Context key which is used to store
// the brokers list passed to the server.go and share it among all contexts
const SEED_CONTEXT_KEY = "seeds"

// Comma separated Kafka broker list, split and stored into string array
type Seeds []string

// Request holds the JSON request to post data to a topic
type Request struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Partition string `json:"partition,omitempty"`
	Topic     string `json:"topic"`
}

// Response holds the JSON response to return to consumer
// Just a wrapper for underlying kgo.Record and kgo.FetchError
type Response struct {
	Key       string           `json:"key"`
	Value     string           `json:"value"`
	Partition string           `json:"partition"`
	Offset    string           `json:"offset"`
	Timestamp string           `json:"timestamp"`
	Errors    []kgo.FetchError `json:"errors,omitempty"`
}

// Query handles the query parameter
type Query struct {
	Topic string `query:"topic"`
	Group string `query:"group"`
}

var sugar *zap.SugaredLogger

func init() {
	logger, _ := zap.NewProduction()
	logger.Sync()
	sugar = logger.Sugar()
}

// SeedsContext adds the broker list to echo.Context and make it available across the context
func (s Seeds) SeedsContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sugar.Infow("Kafka client configuration", "brokers", s)
			c.Set(SEED_CONTEXT_KEY, s)
			return next(c)
		}
	}
}

// NewClient builds the Kafka Client using the provided options
func (s Seeds) NewClient(opts ...kgo.Opt) (*kgo.Client, error) {
	opts = append(opts, kgo.SeedBrokers(s...))
	return kgo.NewClient(opts...)
}
