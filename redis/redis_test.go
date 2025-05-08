package redis

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestRedisSingleClient(t *testing.T) {
	url := "redis://localhost:7001/?dial_timeout=3&&read_timeout=6s&max_retries=2"
	c, err := NewSingleClient(url)
	if err != nil {
		log.Printf("new single client: %v", err)
		return
	}
	defer c.Close()

	if err = c.Ping(context.TODO()).Err(); err != nil {
		log.Printf("ping error: %v", err)
		return
	}

	log.Printf("success")
}

func TestRedisClusterClient(t *testing.T) {
	clusterUrl := "redis://localhost:7001?dial_timeout=3s&read_timeout=6s&addr=localhost:7002&addr=localhost:7003&addr=localhost:7004&addr=localhost:7005&addr=localhost:7006&conn_max_idle_time=5s&conn_max_lifetime=5s"

	cc, err := NewClusterClient(clusterUrl)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer cc.Close()

	fmt.Println(cc.ClientID(context.TODO()).String())
	fmt.Println(cc.PoolStats())
}

func TestRedisUniversalClient(t *testing.T) {
	// single client mode
	_ = &redis.UniversalOptions{
		Addrs: []string{"localhost:7001"},
	}

	// fail over client mode
	uo := &redis.UniversalOptions{
		Addrs:      []string{"localhost:7001", "localhost:7002"},
		MasterName: "master",
	}

	// cluster client mode
	_ = &redis.UniversalOptions{
		Addrs:      []string{"localhost:7001", "localhost:7002"},
		MasterName: "",
	}

	uc := NewUniversalClient(uo)
	defer uc.Close()

	if err := uc.Ping(context.TODO()).Err(); err != nil {
		log.Printf("ping error: %v", err)
	}
}
