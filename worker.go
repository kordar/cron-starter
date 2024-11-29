package cron_starter

import (
	"github.com/kordar/goresty"
)

type FetchMethod struct {
	Name  string      `json:"name"`
	JobId string      `json:"job_id"`
	Param interface{} `json:"param"`
}

type RemoteWorker struct {
	id    string
	feign *goresty.Feign
}

func (r *RemoteWorker) Id() string {
	return r.id
}

func (r *RemoteWorker) Feign() *goresty.Feign {
	return r.feign
}

func NewRemoteWorker(id string, feign *goresty.Feign) *RemoteWorker {
	return &RemoteWorker{id, feign}
}
