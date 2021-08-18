package utils

import (
	"context"
	"time"

	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/core"
	"github.com/mimuret/golang-iij-dpf/pkg/apiutils"
)

func Client(logger api.Logger) (*api.Client, error) {
	c, err := GetContexts()
	if err != nil {
		return nil, err
	}
	return api.NewClient(c.Token, c.Endpoint, logger), nil
}

func Wait(c *api.Client, jobId string, timeout time.Duration) (*core.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return apiutils.WaitJob(ctx, c, jobId, time.Second)
}
