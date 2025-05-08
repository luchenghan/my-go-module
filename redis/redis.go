package redis

import (
	"github.com/redis/go-redis/v9"
)

func NewSingleClient(url string) (*redis.Client, error) {
	o, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(o), nil
}

func NewClusterClient(url string) (*redis.ClusterClient, error) {
	o, err := redis.ParseClusterURL(url)
	if err != nil {
		return nil, err
	}

	return redis.NewClusterClient(o), nil
}

func NewUniversalClient(o *redis.UniversalOptions) redis.UniversalClient {
	return redis.NewUniversalClient(o)
}
