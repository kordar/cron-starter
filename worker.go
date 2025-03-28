package cron_starter

import (
	"github.com/go-resty/resty/v2"
	goframeworkcron "github.com/kordar/goframework-cron"
	goframeworkresty "github.com/kordar/goframework-resty"
	logger "github.com/kordar/gologger"
	"github.com/kordar/goresty"
	"github.com/spf13/cast"
	"time"
)

func loadWorker(name string, options *StarterOptions) {

	if options.NodeType != "worker" {
		return
	}

	feign := getFeign(name, options)
	logger.Infof("===============%s", options.Id)
	goframeworkcron.AddJob(options.Id, NewWorkerHeartSchedule(feign, options))
	//goframeworkcron.AddJob(id, NewWorkerHeartSchedule(id, cfg["remote_worker_heart_spec"]))

}

func getFeign(name string, options *StarterOptions) *goresty.Feign {
	if options.NodeId == "" {
		logger.Fatalf("[%s] you must configure the parameter \"node_id\"", name)
	}

	if options.WorkerFeignHost == "" {
		logger.Fatalf("[%s] you must configure the parameter \"worker_feign_host\"", name)
	}

	if options.WorkerFeign != "" && goframeworkresty.HasFeignInstance(options.WorkerFeign) {
		return goframeworkresty.GetFeignClient(options.WorkerFeign)
	}

	feign := goresty.NewFeign(nil)
	feign.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		if options.WorkerFeignTrace == "enable" {
			request.EnableTrace()
		}
		return nil
	})

	feign.Options(func(client *resty.Client) {
		client.SetBaseURL(options.WorkerFeignHost)
		if options.WorkerFeignDebug == "enable" {
			client.SetDebug(true)
		}
		if options.WorkerFeignTimeout != 0 {
			remoteTimeout := cast.ToDuration(options.WorkerFeignTimeout)
			client.SetTimeout(time.Second * remoteTimeout)
		}
		if options.WorkerFeignRetryCount != 0 {
			client.SetRetryCount(options.WorkerFeignRetryCount)
		}
		if options.WorkerFeignRetryWaitTime != 0 {
			remoteRetryWaitTime := cast.ToDuration(options.WorkerFeignRetryWaitTime)
			client.SetRetryWaitTime(time.Second * remoteRetryWaitTime)
		}
	})

	feign.OnError(func(request *resty.Request, err error) {
		logger.Errorf("[%s] request err = %+v", name, err)
	})

	return feign
}

type FetchMethod struct {
	Name  string      `json:"name"`
	JobId string      `json:"job_id"`
	Param interface{} `json:"param"`
}
