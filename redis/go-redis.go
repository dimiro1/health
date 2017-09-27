package redis

import (
	"errors"
	"strings"

	"github.com/go-redis/redis"
)

type goRedis struct {
	client *redis.Client
}

// NewGoRedis returns a Redis implementation that gets the version of
// github.com/go-redis/redis client
func NewGoRedis(client *redis.Client) Redis {
	return goRedis{
		client: client,
	}
}

func (r goRedis) GetVersion() (string, error) {
	info, err := r.client.Info("server").Result()
	if err != nil {
		return "", err
	}

	version, err := getVersion(info)
	if err != nil {
		return "", nil
	}
	return version, nil
}

func getVersion(info string) (string, error) {
	split := strings.Split(info, "\r\n")

	for _, line := range split {
		values := strings.Split(line, ":")
		if len(values) == 2 {
			if values[0] == "redis_version" {
				return values[1], nil
			}
		}
	}

	return "", errors.New("could not find Redis version")
}
