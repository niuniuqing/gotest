package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/go-redis/redis"
	"os"
	"startchartsDemo/config"
)

func main() {

	log.SetHandler(text.New(os.Stderr))
	log.SetLevel(log.DebugLevel)
	config := config.Get()
	//返回一个*entry
	ctx := log.WithField("listen", config.Listen)
	options, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		log.WithError(err).Fatal("invalid redis_url")
	}
	//返回一个*click
	redis := redis.NewClient(options)
	cache = cache.New(redis)

}
