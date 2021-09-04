package utils

import (
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/core"
	"github.com/spf13/viper"
	"go.uber.org/multierr"
)

var _ api.Spec = &CommandResults{}

type CommandResults struct {
	Items []CommandResult
	Err   error
}

func (a *CommandResults) DeepCopyObject() api.Object {
	c := &CommandResults{}
	c.Items = append(c.Items, a.Items...)
	return c
}
func (a *CommandResults) Add(reqId string, err error) {
	if err != nil {
		a.Err = multierr.Append(a.Err, err)
	}
	a.Items = append(a.Items, CommandResult{
		RequestId: reqId,
		Err:       err,
	})
}
func (a *CommandResults) GetName() string                             { return "" }
func (a *CommandResults) GetGroup() string                            { return "" }
func (a *CommandResults) GetPathMethod(_ api.Action) (string, string) { return "", "" }
func (a *CommandResults) WaitJob(cl api.ClientInterface, v *viper.Viper) {
	if !v.GetBool("wait") {
		return
	}
	for i, result := range a.Items {
		if result.Err == nil {
			job, err := Wait(cl, result.RequestId, v.GetDuration("wait-timeout"))
			a.Items[i].Job = job
			a.Items[i].Err = err
			if err != nil {
				a.Err = multierr.Append(a.Err, err)
			}
		}
	}
}

type CommandResult struct {
	RequestId string
	Err       error
	Job       *core.Job
}
