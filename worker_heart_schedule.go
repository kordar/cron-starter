package cron_starter

import (
	"github.com/kordar/gocron"
	goframeworkcron "github.com/kordar/goframework-cron"
	logger "github.com/kordar/gologger"
	"github.com/kordar/goresty"
	"strings"
)

type WorkerHeartSchedule struct {
	feign   *goresty.Feign
	options *StarterOptions
	*gocron.BaseSchedule
}

func NewWorkerHeartSchedule(feign *goresty.Feign, options *StarterOptions) *WorkerHeartSchedule {
	return &WorkerHeartSchedule{feign, options, &gocron.BaseSchedule{}}
}

func (h *WorkerHeartSchedule) GetId() string {
	return "#worker-heartbeat"
}

func (h *WorkerHeartSchedule) GetSpec() string {
	if h.options.WorkerHeartbeatSpec != "" {
		return h.options.WorkerHeartbeatSpec
	} else {
		return h.BaseSchedule.GetSpec()
	}
}

func (h *WorkerHeartSchedule) Execute() {
	items := goframeworkcron.StateJob(h.options.Id)
	fetchMethod := FetchMethod{}
	_, err := h.feign.Request().SetBody(items).SetResult(&fetchMethod).Post(h.options.WorkerHeartbeatUrl)
	if err != nil {
		logger.Warnf("[%s-%s] request %s err: %v", h.options.Id, h.GetId(), h.options.WorkerHeartbeatUrl, err)
		return
	}

	logger.Infof("[%s-%s] response data is %+v", h.options.Id, h.GetId(), fetchMethod)

	if strings.HasPrefix(fetchMethod.Name, "#") {
		return
	}

	if fetchMethod.Name == "stop-job" {
		goframeworkcron.RemoveJob(h.options.Id, fetchMethod.JobId)
	}

	if fetchMethod.Name == "reload-job" {
		goframeworkcron.ReloadJob(h.options.Id, fetchMethod.JobId)
	}

}
