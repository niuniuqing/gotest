package cache

import (
	rediscache "github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vmihailenco/msgpack"
)

var cacheGets = prometheus.NewCounter(
	prometheus.CounterOpts{
		Namespace: "starchartsDemo",
		Subsystem: "cache",
		Name:      "gets_total",
		Help:      "Total number of successful cache gets",
	},
)

var cachePuts = prometheus.NewCounter(
	prometheus.CounterOpts{
		Namespace: "starchartsDemo",
		Subsystem: "cache",
		Name:      "puts_total",
		Help:      "Total number of successful cache puts",
	},
)

var cacheDeletes = prometheus.NewCounter(
	prometheus.CounterOpts{
		Namespace: "starcharts",
		Subsystem: "cache",
		Name:      "deletes_total",
		Help:      "Total number of successful cache deletes",
	},
)

func init() {
	prometheus.MustRegister(cacheGets, cachePuts)
}

type Redis struct {
	redis *redis.Client
	codec *rediscache.Codec
}

func New(redis *redis.Client) *Redis {
	codec := &rediscache.Codec{
		Redis: redis,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
	return &Redis{
		redis: redis,
		codec: codec,
	}
}

func (c *Redis) Close() error {
	return c.redis.Close()
}

func (c *Redis) Get(key string, result interface{}) error {
	if err := c.codec.Get(key, result); err != nil {
		return err
	}
	cacheGets.Inc()
	return nil
}
