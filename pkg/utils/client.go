package utils

import (
	"time"

	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/core"
	"github.com/mimuret/golang-iij-dpf/pkg/apiutils"
)

func Client(logger api.Logger) (*api.Client, error) {
	c, err := GetConfig()
	if err != nil {
		return nil, err
	}
	return api.NewClient(c.Token, c.Endpoint, logger), nil
}

func Wait(c *api.Client, jobId string, timeout time.Duration) (*core.Job, error) {
	c.SetWatchTimeout(timeout)
	return apiutils.WaitJob(c, jobId)
}
