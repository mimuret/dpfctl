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
	"github.com/spf13/cobra"
)

var CreateOption = &createOpt{}

type createOpt struct {
	utils.CommonParam
	Filename string
}

func newCreateCmd() *cobra.Command {
	cmd := utils.NewCommand("create resource_name -f filename", api.ActionCreate, runCreate)
	cmd.Args = cobra.MinimumNArgs(1)
	utils.ChangeCmd(cmd, &CreateOption.CommonParam)
	cmd.PersistentFlags().StringVarP(&CreateOption.Filename, "filename", "f", "", "resource data file")
	cmd.MarkPersistentFlagFilename("filename", "yaml", "yml", "json")
	cmd.RegisterFlagCompletionFunc("filename", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveFilterFileExt
	})
	return cmd
}

func runCreate(cmd *cobra.Command, args []string, spec apis.Spec) error {
	cobra.CheckErr(utils.ReadSpec(spec, CreateOption.Filename))
	cl, err := utils.NewClient(log)
	if err != nil {
		return err
	}
	if CreateOption.DryRun {
		fmt.Printf("[DryRun] create %v", spec)
		return nil
	}
	reqId, err := cl.Create(spec, nil)
	if err != nil {
		return err
	}
	if CreateOption.Wait {
		job, err := utils.Wait(cl, reqId, CreateOption.WaitTimeout)
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
