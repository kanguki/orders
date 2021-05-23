package cache

// var redis
import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var REDIS_HOST = os.Getenv("REDIS_HOST")
var REDIS_PORT = os.Getenv("REDIS_PORT")

var ctx = context.Background()
var Redis *redis.Client = nil

func GetRedis() *redis.Client {
	addr := fmt.Sprintf("%s:%s",REDIS_HOST,REDIS_PORT)
	if Redis == nil {
		redis := redis.NewClient(&redis.Options{
			Addr: addr,
		})
		_, err := redis.Ping(ctx).Result()
		if err != nil {
			log.Fatal("Error connecting redis: ", err)
		}
		Redis = redis
		log.Println("Cache is ready")
	}
	return Redis
}