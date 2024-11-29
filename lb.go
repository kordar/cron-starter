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
	Can(value string) bool
	Load([]string)
}

func StartRedisRegistry(name string, prefix string, nodeStr string, timeoutSeconds int, heartbeatSeconds int, weight string, channel string, lb LB) {

	if !goframeworkredis.HasRedisInstance(name) {
		logger.Fatalf("To start the redis registry, you must running the instance '%s' for redis.", name)
	}

	if timeoutSeconds <= 0 {
		timeoutSeconds = 30
	}

	if heartbeatSeconds <= 0 {
		heartbeatSeconds = 300
	}

	if prefix == "" {
		logger.Fatal("To start the redis registry, you must configure the parameter 'lb_redis_prefix'.")
	}

	if channel == "" {
		logger.Fatal("To start the redis registry, you must configure the parameter 'lb_redis_channel'.")
	}

	redisClient := goframeworkredis.GetRedisClient(name)
	var redisnoderegistry registry.Registry = registry_goredis.NewRedisNodeRegistry(redisClient, &registry_goredis.RedisNodeRegistryOptions{
		Prefix:  prefix,
		Node:    nodeStr,
		Timeout: time.Second * time.Duration(timeoutSeconds),
		Channel: channel,
		Weight:  weight,
		Reload: func(value []string, channel string) {
			lb.Load(value)
		},
		Heartbeat: time.Second * time.Duration(heartbeatSeconds),
	})
	redisnoderegistry.Listener()
	_ = redisnoderegistry.Register()
}
