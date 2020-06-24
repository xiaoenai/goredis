package goredis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/henrylee2cn/cfgo"
)

type (
	// Config redis (cluster) client config
	Config struct {
		// redis deploy type, [single, cluster]
		DeployType string `yaml:"deploy_type"`
		// only for single node config, valid when DeployType=single.
		ForSingle SingleConfig `yaml:"for_single"`
		// only for cluster config, valid when DeployType=cluster.
		ForCluster ClusterConfig `yaml:"for_cluster"`

		// An optional password. Must match the password specified in the
		// requirepass server configuration option.
		Password string `yaml:"password,omitempty"`

		// The maximum number of retries before giving up.
		// Default is to not retry failed commands.
		MaxRetries int `yaml:"max_retries,omitempty"`

		// Dial timeout for establishing new connections.
		// Default is 5 seconds.
		DialTimeout int64 `yaml:"dial_timeout,omitempty"`
		// Timeout for socket reads. If reached, commands will fail
		// with a timeout instead of blocking.
		// Default is 3 seconds.
		ReadTimeout int64 `yaml:"read_timeout,omitempty"`
		// Timeout for socket writes. If reached, commands will fail
		// with a timeout instead of blocking.
		// Default is ReadTimeout.
		WriteTimeout int64 `yaml:"write_timeout,omitempty"`

		// PoolSizePerNode applies per cluster node and not for the whole cluster.
		// Maximum number of socket connections.
		// Default is 10 connections per every CPU as reported by runtime.NumCPU.
		PoolSizePerNode int `yaml:"pool_size_per_node"`
		// Amount of time client waits for connection if all connections
		// are busy before returning an error.
		// Default is ReadTimeout + 1 second.
		PoolTimeout int64 `yaml:"pool_timeout,omitempty"`
		// Amount of time after which client closes idle connections.
		// Should be less than server's timeout.
		// Default is 300 seconds.
		IdleTimeout int64 `yaml:"idle_timeout"`
		// Frequency of idle checks.
		// Default is 60 seconds.
		// When minus value is set, then idle check is disabled.
		IdleCheckFrequency int64 `yaml:"idle_check_frequency,omitempty"`

		// Enables read only queries on slave nodes.
		// Only for cluster.
		ReadOnly bool `yaml:"read_only,omitempty"`

		init bool
	}

	// SingleConfig redis single node client config.
	SingleConfig struct {
		// host:port address.
		Addr string `yaml:"addr"`

		// Maximum backoff between each retry.
		// Default is 512 seconds; -1 disables backoff.
		MaxRetryBackoff int64 `yaml:"max_retry_backoff,omitempty"`
	}

	// ClusterConfig redis cluster client config.
	ClusterConfig struct {
		// A seed list of host:port addresses of cluster nodes.
		Addrs []string `yaml:"addrs"`

		// The maximum number of retries before giving up. Command is retried
		// on network errors and MOVED/ASK redirects.
		// Default is 16.
		MaxRedirects int `yaml:"max_redirects,omitempty"`

		// Enables routing read-only queries to the closest master or slave node.
		RouteByLatency bool `yaml:"route_by_latency,omitempty"`
	}
)

// deploy types
const (
	TypeSingle  = "single"
	TypeCluster = "cluster"
)

// Reload reloads config.
func (cfg *Config) Reload(bind cfgo.BindFunc) error {
	if cfg.init {
		return nil
	}
	err := bind()
	if err != nil {
		return err
	}
	cfg.init = true
	if cfg.DeployType != TypeSingle && cfg.DeployType != TypeCluster {
		return fmt.Errorf("redis config: deploy_type optional enumeration list: %s, %s", TypeSingle, TypeCluster)
	}
	return nil
}

// ReadConfig read config from specified yaml section.
func ReadConfig(configSection string) (*Config, error) {
	var cfg = NewConfig()
	var err error
	if cfgo.IsReg(configSection) {
		err = cfgo.BindSection(configSection, cfg)
	} else {
		err = cfgo.Reg(configSection, cfg)
	}
	return cfg, err
}

// NewConfig creates a default config.
func NewConfig() *Config {
	return &Config{
		DeployType: TypeSingle,
		ForSingle:  SingleConfig{Addr: "127.0.0.1:6379"},
		// ForCluster: ClusterConfig{Addrs: []string{"127.0.0.1:6379"}},
	}
}