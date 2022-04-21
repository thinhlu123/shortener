package cache

import (
	"context"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/thinhlu123/shortener/config"
	"time"
)

type Cache struct {
	redisCache *cache.Cache
}

func (c *Cache) Init() {
	if c.redisCache != nil {
		ring := redis.NewRing(&redis.RingOptions{
			Addrs: map[string]string{
				config.Conf.Redis.RedisAddr: config.Conf.Redis.RedisPort,
			},
			Username:    config.Conf.Redis.RedisUser,
			Password:    config.Conf.Redis.RedisPassword,
			DB:          config.Conf.Redis.DB,
			PoolSize:    config.Conf.Redis.PoolSize,
			PoolTimeout: time.Second * time.Duration(config.Conf.Redis.PoolTimeout),
		})

		c.redisCache = cache.New(&cache.Options{
			Redis:      ring,
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		})
	}
}

func (c *Cache) Set(key string, obj interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	if err := c.redisCache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: obj,
		TTL:   time.Hour,
	}); err != nil {
		return err
	}

	return nil
}

func (c *Cache) Get(key string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	var rs interface{}
	if err := c.redisCache.Get(ctx, key, &rs); err != nil {
		return nil, err
	}

	return rs, nil
}
