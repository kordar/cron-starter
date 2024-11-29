package cron_starter

import (
	"encoding/json"
	"github.com/kordar/registry"
	"github.com/spf13/cast"
)

type LBNodeWithRedis struct {
	nodeId   string
	hashring registry.HashringRegistry
}

func NewLBNodeWithRedis(nodeId string, spots int) *LBNodeWithRedis {
	return &LBNodeWithRedis{
		nodeId:   nodeId,
		hashring: registry.NewHashringRegistry(spots, nodeId),
	}
}

func (n *LBNodeWithRedis) Can(key string) bool {
	return n.hashring.Can(key)
}

func (n *LBNodeWithRedis) Load(value []string) {
	n.hashring.Load(value, func(v interface{}) (string, int, []string) {
		str := cast.ToString(v)
		node := LBNode{}
		if err := json.Unmarshal([]byte(str), &node); err == nil {
			return node.Id, node.Weight, nil
		} else {
			return "", 0, nil
		}
	})
}
