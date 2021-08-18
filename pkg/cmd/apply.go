/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis"
	"github.com/mimuret/golang-iij-dpf/pkg/schema"
	"github.com/spf13/cobra"
)

var ApplyOption = &applyOpt{}

type applyOpt struct {
	utils.CommonParam
	Filename string
}

func newCmdApply() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "apply -f filename",
		Args: cobra.NoArgs,
		RunE: runApply,
	}
	utils.ChangeCmd(cmd, &ApplyOption.CommonParam)
	cmd.PersistentFlags().StringVarP(&ApplyOption.Filename, "filename", "f", "", "resource data file")
	cmd.MarkPersistentFlagFilename("filename", "yaml", "yml", "json")
	cmd.RegisterFlagCompletionFunc("filename", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveFilterFileExt
	})
	return cmd
}

func runApply(cmd *cobra.Command, _ []string) error {
	bs, err := utils.ReadFile(ApplyOption.Filename)
	spec, err := schema.SchemaSet.Parse(bs)
	if err != nil {
		return err
	}
	cl, err := utils.Client(log)
	if err != nil {
		return err
	}
	apisSpec, ok := spec.(apis.Spec)
	if !ok {
		return fmt.Errorf("not support file, this spec is not apis.Spec")
	}
	var (
		create bool
		update bool
		exist  bool
		reqId  string
	)
	if method, _ := apisSpec.GetPathMethod(api.ActionCreate); method != "" {
		create = true
	}
	if method, _ := apisSpec.GetPathMethod(api.ActionUpdate); method != "" {
		update = true
	}
	if !create && !update {
		return fmt.Errorf("this spec is not support create and update")
	}

	_, err = cl.Read(apisSpec)
	if api.IsNotFound(err) {
		exist = true
	}
	if exist && !update {
		return fmt.Errorf("this spec is not support update")
	}
	if !exist && create {
		return fmt.Errorf("this spec is not support create")
	}
	if exist {
		if ApplyOption.DryRun {
			fmt.Printf("[DryRun] update resource %v", apisSpec)
		} else {
			reqId, err = cl.Update(apisSpec, nil)
		}
	} else {
		if ApplyOption.DryRun {
			fmt.Printf("[DryRun] create resource %v", apisSpec)
		} else {
			reqId, err = cl.Create(apisSpec, nil)
		}
	}
	if err != nil {
		return fmt.Errorf("failed to request JobId: %s err: %w", reqId, err)
	}
	if ApplyOption.Wait {
		job, err := utils.Wait(cl, reqId, ApplyOption.WaitTimeout)
		if err != nil {
			return fmt.Errorf("failed to request JobId: %s err: %w", reqId, err)
		}
		log.Infof("success request\n")
		log.Infof("JobId: %s\n", reqId)
		log.Infof("ResourceURL: %s\n", job.ResourceUrl)
	} else {
		log.Infof("success request JobId: %s\n", reqId)
	}
	return nil
}
