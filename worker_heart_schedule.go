package cron_starter

import (
	"github.com/kordar/gocron"
	goframeworkcron "github.com/kordar/goframework-cron"
	logger "github.com/kordar/gologger"
	"strings"
)

type WorkerHeartSchedule struct {
	id   string
	spec string
	url  string
	*gocron.BaseSchedule
}

func NewWorkerHeartSchedule(id string, spec string, url string) *WorkerHeartSchedule {
	return &WorkerHeartSchedule{id, spec, url, &gocron.BaseSchedule{}}
}

func (h *WorkerHeartSchedule) GetId() string {
	return "::worker_heart"
}

func (h *WorkerHeartSchedule) GetSpec() string {
	if h.spec != "" {
		return h.spec
	} else {
		return h.BaseSchedule.GetSpec()
	}
}

func (h *WorkerHeartSchedule) Execute() {
	remoteWorker := RemoteWorkerHandle[h.id]
	if remoteWorker == nil {
		logger.Warnf("[%s-%s] no valid object found for RemoteWorker", h.id, h.GetId())
		return
	}
	items := goframeworkcron.StateJob(h.id)
	fetchMethod := FetchMethod{}
	_, err := remoteWorker.Feign().Request().SetBody(items).SetResult(&fetchMethod).Post(h.url)
	if err != nil {
		logger.Warnf("[%s-%s] request %s err: %v", h.id, h.GetId(), h.url, err)
		return
	}

	logger.Infof("[%s-%s] response data is %+v", h.id, h.GetId(), fetchMethod)

	if strings.HasPrefix(fetchMethod.Name, "::") {
		return
	}

	if fetchMethod.Name == "stop-job" {
		goframeworkcron.RemoveJob(h.id, fetchMethod.JobId)
	}

	if fetchMethod.Name == "reload-job" {
		goframeworkcron.ReloadJob(h.id, fetchMethod.JobId)
	}

}
