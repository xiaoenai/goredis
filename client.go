package goredis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

// Client redis (cluster) client.
type (
	Client struct {
		cfg *Config
		Cmdable
	}
	Cmdable interface {
		redis.Cmdable
		Subscribe(channels ...string) *redis.PubSub
	}
	// Alias
	PubSub             = redis.PubSub
	Message            = redis.Message
	GeoLocation        = redis.GeoLocation
	GeoRadiusQuery     = redis.GeoRadiusQuery
	ZRangeBy           = redis.ZRangeBy
	Z                  = redis.Z
	Pipeliner          = redis.Pipeliner
	RedisCmdable       = redis.Cmdable
	SliceCmd           = redis.SliceCmd
	StatusCmd          = redis.StatusCmd
	Cmder              = redis.Cmder
	IntCmd             = redis.IntCmd
	DurationCmd        = redis.DurationCmd
	BoolCmd            = redis.BoolCmd
	StringCmd          = redis.StringCmd
	FloatCmd           = redis.FloatCmd
	StringSliceCmd     = redis.StringSliceCmd
	BoolSliceCmd       = redis.BoolSliceCmd
	StringStringMapCmd = redis.StringStringMapCmd
	StringIntMapCmd    = redis.StringIntMapCmd
	ZSliceCmd          = redis.ZSliceCmd
	ScanCmd            = redis.ScanCmd
	ClusterSlotsCmd    = redis.ClusterSlotsCmd
)

// NewClient creates a redis(cluster) client from yaml config, and pings the client.
func NewClient(cfg *Config) (*Client, error) {
	var c = &Client{
		cfg: cfg,
	}
	switch cfg.DeployType {
	case TypeSingle:
		c.Cmdable = redis.NewClient(&redis.Options{
			Addr:               cfg.ForSingle.Addr,
			Password:           cfg.Password,
			MaxRetries:         cfg.MaxRetries,
			MaxRetryBackoff:    time.Duration(cfg.ForSingle.MaxRetryBackoff) * time.Second,
			DialTimeout:        time.Duration(cfg.DialTimeout) * time.Second,
			ReadTimeout:        time.Duration(cfg.ReadTimeout) * time.Second,
			WriteTimeout:       time.Duration(cfg.WriteTimeout) * time.Second,
			PoolSize:           cfg.PoolSizePerNode,
			PoolTimeout:        time.Duration(cfg.PoolTimeout) * time.Second,
			IdleTimeout:        time.Duration(cfg.IdleTimeout) * time.Second,
			IdleCheckFrequency: time.Duration(cfg.IdleCheckFrequency) * time.Second,
		})

	case TypeCluster:
		c.Cmdable = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:              cfg.ForCluster.Addrs,
			MaxRedirects:       cfg.ForCluster.MaxRedirects,
			ReadOnly:           cfg.ReadOnly,
			RouteByLatency:     cfg.ForCluster.RouteByLatency,
			Password:           cfg.Password,
			MaxRetries:         cfg.MaxRetries,
			DialTimeout:        time.Duration(cfg.DialTimeout) * time.Second,
			ReadTimeout:        time.Duration(cfg.ReadTimeout) * time.Second,
			WriteTimeout:       time.Duration(cfg.WriteTimeout) * time.Second,
			PoolSize:           cfg.PoolSizePerNode,
			PoolTimeout:        time.Duration(cfg.PoolTimeout) * time.Second,
			IdleTimeout:        time.Duration(cfg.IdleTimeout) * time.Second,
			IdleCheckFrequency: time.Duration(cfg.IdleCheckFrequency) * time.Second,
		})

	default:
		return nil, fmt.Errorf("redis.Config.DeployType: optional enumeration list: %s, %s", TypeSingle, TypeCluster)
	}

	if _, err := c.Ping().Result(); err != nil {
		return nil, err
	}
	return c, nil
}

// Config returns config.
func (c *Client) Config() *Config {
	return c.cfg
}

// IsCluster returns whether it is a cluster.
func (c *Client) IsCluster() bool {
	return c.cfg.DeployType == TypeCluster
}

// ToSingle tries to convert it to *redis.Client.
func (c *Client) ToSingle() (*redis.Client, bool) {
	cli, ok := c.Cmdable.(*redis.Client)
	return cli, ok
}

// ToCluster tries to convert it to *redis.ClusterClient.
func (c *Client) ToCluster() (*redis.ClusterClient, bool) {
	clu, ok := c.Cmdable.(*redis.ClusterClient)
	return clu, ok
}

// LockCallback 使用分布式锁执行回调函数
// 注意：每10毫秒尝试1次上锁，且上锁后默认锁定1分钟
func (c *Client) LockCallback(lockKey string, callback func(), maxLock ...time.Duration) error {
	var d = time.Minute
	if len(maxLock) > 0 {
		d = maxLock[0]
	}
	// lock
	for lockOk, err := c.SetNX(lockKey, "", d).Result(); !lockOk; lockOk, err = c.SetNX(lockKey, "", d).Result() {
		if err != nil && !IsRedisNil(err) {
			return err
		}
		time.Sleep(time.Millisecond * 10)
	}
	// unlock
	defer c.Del(lockKey)
	// do
	callback()
	return nil
}

// Redis nil reply, .e.g. when key does not exist.
const Nil = redis.Nil

// IsRedisNil Is the redis nil reply? .e.g. when key does not exist.
func IsRedisNil(err error) bool {
	return redis.Nil == err
}
