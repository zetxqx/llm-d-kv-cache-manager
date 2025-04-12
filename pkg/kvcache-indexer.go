package pkg

import (
	"github.com/neuralmagic/distributed-kv-cache/pkg/utils"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"

	"strings"
)

// KVCacheIndexer defines the interactions with the KVCache indexing backend.
type KVCacheIndexer interface {
	// GetPodsForKeys retrieves the pods for the given keys.
	// It returns a map where the keys are CacheEngineKey and the values are slices of pod names.
	GetPodsForKeys(ctx context.Context, keys []CacheEngineKey) (map[CacheEngineKey][]string, error)
}

type RedisKVCacheIndexer struct {
	// RedisClient is the Redis client used for communication.
	RedisClient *redis.Client
}

// NewRedisKVCacheIndexer creates a new RedisKVCacheIndexer instance.
func NewRedisKVCacheIndexer(redisClient *redis.Client) *RedisKVCacheIndexer {
	return &RedisKVCacheIndexer{
		RedisClient: redisClient,
	}
}

// GetPodsForKeys retrieves the pods for the given keys.
func (r *RedisKVCacheIndexer) GetPodsForKeys(ctx context.Context, keys []CacheEngineKey) (map[CacheEngineKey][]string, error) {
	pods := make(map[CacheEngineKey][]string)

	redisKeys := utils.SliceMap(keys, func(key CacheEngineKey) string {
		return key.String()
	})
	// use redis.MGet to get all keys at once
	values, err := r.RedisClient.MGet(ctx, redisKeys...).Result()
	if err != nil {
		return nil, err
	}

	for i, value := range values { // values are "podIP:port", we only need podIP
		if value == "" {
			continue
		}

		valueStr, ok := value.(string)
		if !ok {
			continue
		}

		parts := strings.Split(valueStr, ":")
		if len(parts) != 2 {
			continue
		}

		pods[keys[i]] = append(pods[keys[i]], parts[0])
	}

	return pods, nil
}
