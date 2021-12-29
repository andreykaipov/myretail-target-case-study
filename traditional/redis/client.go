package redis

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var URI string
var Client *redis.Client

func init() {
	URI = os.Getenv("REDIS_URI")

	if URI == "" {
		log.Println("REDIS_URI was not provided, so caching won't be possible.")
		return
	}

	opt, err := redis.ParseURL(URI)
	if err != nil {
		log.Fatalf("Failed parsing Redis URI: %s", err)
	}

	Client = redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	})
}
