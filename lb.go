package cron_starter

import (
	goframeworkredis "github.com/kordar/goframework-goredis"
	logger "github.com/kordar/gologger"
	"github.com/kordar/registry"
	"github.com/kordar/registry-goredis"
	"time"
)

type LBNode struct {
	Id          string    `json:"node"`
	RefreshTime time.Time `json:"refresh_time"`
	Status      string    `json:"status"`
	Weight      int       `json:"weight"`
}

type LB interface {
	Can(key string) bool
	Load(value []string)
}

func loadLb(name string, options *StarterOptions) LB {
	if options.Lb == "" {
		return nil
	}

	if options.NodeId == "" {
		logger.Fatalf("[%s] to enable load balancing, you must configure the parameter \"node_id\"", name)
	}

	if options.Lb == "redis" {
		lbNodeWithRedis := NewLBNodeWithRedis(options.NodeId, options.HashringSpots)
		StartRedisRegistry(options, lbNodeWithRedis)
		return lbNodeWithRedis
	}

	return nil
}

func StartRedisRegistry(options *StarterOptions, lb LB) {

	if !goframeworkredis.HasRedisInstance(options.LbRedis) {
		logger.Fatalf("To start the redis registry, you must running the instance '%s' for redis.", options.LbRedis)
	}

	if options.LbRedisPrefix == "" {
		logger.Fatal("To start the redis registry, you must configure the parameter 'lb_redis_prefix'.")
	}

	if options.LbRedisChannel == "" {
		logger.Fatal("To start the redis registry, you must configure the parameter 'lb_redis_channel'.")
	}

	redisClient := goframeworkredis.GetRedisClient(options.LbRedis)
	var redisnoderegistry registry.Registry = registry_goredis.NewRedisNodeRegistry(redisClient, &registry_goredis.RedisNodeRegistryOptions{
		Prefix:  options.LbRedisPrefix,
		Node:    options.NodeId,
		Timeout: time.Second * time.Duration(options.LbRedisTimeout),
		Channel: options.LbRedisChannel,
		Weight:  options.LbRedisWeight,
		Reload: func(value []string, channel string) {
			lb.Load(value)
		},
		Heartbeat: time.Second * time.Duration(options.LbRedisHeartbeat),
	})
	redisnoderegistry.Listener()
	_ = redisnoderegistry.Register()
}
