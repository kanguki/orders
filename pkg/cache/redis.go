package cache

// var redis
import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var REDIS_HOST = os.Getenv("REDIS_HOST")
var REDIS_PORT = os.Getenv("REDIS_PORT")

var doOnce sync.Once
var ctx = context.Background()
var Redis *redis.Client = nil

type Cache struct {
	client *redis.Client
}

func GetCache() *Cache {
	return &Cache{
		client: connectRedis(),
	}
}

func connectRedis() *redis.Client {
	addr := fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT)
	doOnce.Do(func() {
		redis := redis.NewClient(&redis.Options{
			Addr: addr,
		})
		_, err := redis.Ping(ctx).Result()
		if err != nil {
			log.Fatal("Error connecting redis: ", err)
		}
		Redis = redis
		log.Println("Cache is ready")
	})
	return Redis
}
func (cache *Cache) SetEX(key string, value interface{}, seconds int16) interface{} {
	return cache.client.SetEX(ctx, key, value, time.Duration(seconds)*time.Second)
}
func (cache *Cache) Set(key string, value interface{}) interface{} {
	return cache.client.Set(ctx, key, value, 0)
}
func (cache *Cache) Get(key string) interface{} {
	val, err := cache.client.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("key %v does not exist", key)
		return nil
	} else if err != nil {
		log.Printf("error querying redis: %v", err)
	}
	return val
}
