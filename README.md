# goredis
Golang Redis tools, Build based on `github.com/go-redis/redis/v7` v7.0.0-beta.4

## Example

```
package main

import (
	"log"
	"time"

	"github.com/xiaoenai/goredis/v2"
)

func main() {
	cfg, err := redis.ReadConfig("test_redis")
	if err != nil {
		log.Fatalf("redis.ReadConfig(\"test_redis\"): %v", err)
	}
	c, err := redis.NewClient(cfg)
	if err != nil {
		log.Fatalf("redis.NewClient(\"test_redis\"): %v", err)
	}

	m := redis.NewModule("test", "v1.0")

	s, err := c.Set(m.Key("a_key"), "a_value", time.Second).Result()
	if err != nil {
		log.Fatalf("c.Set().Result() error: %v", err)
	}
	log.Printf("c.Set().Result() result: %s", s)

	s, err = c.Get(m.Key("a_key")).Result()
	if err != nil {
		log.Fatalf("c.Get().Result() error: %v", err)
	}
	log.Printf("c.Get().Result() result: %s", s)
	time.Sleep(2 * time.Second)

	s, err = c.Get(m.Key("a_key")).Result()
	if err == nil {
		log.Fatalf("[after 2s] c.Get().Result() result: %s", s)
	}
	log.Printf("[after 2s] c.Get().Result() error: %s", err)
}
```

## API doc

[http://godoc.org/gopkg.in/go-redis/redis.v7](http://godoc.org/gopkg.in/go-redis/redis.v7)
