package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/core"
	"github.com/mimuret/golang-iij-dpf/pkg/apiutils"
)

var (
	NewClient  = NewClientDefault
	ErrTimeout = fmt.Errorf("timeout")
)

func NewClientDefault(logger api.Logger) (api.ClientInterface, error) {
	c, err := GetContexts()
	if err != nil {
		return nil, err
	}
	return api.NewClient(c.Token, c.Endpoint, logger), nil
}

func Wait(c api.ClientInterface, jobId string, timeout time.Duration) (*core.Job, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	job, err := apiutils.WaitJob(ctx, c, jobId, time.Second)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ErrTimeout
		}
		return job, err
	}
	return job, nil
}
