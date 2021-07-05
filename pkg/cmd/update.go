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
	"os"

	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis"
	"github.com/spf13/cobra"
)

var UpdateOption = &updateOpt{}

type updateOpt struct {
	utils.CommonParam
	Filename string
}

func newUpdateCmd() *cobra.Command {
	cmd := utils.NewCommand("update resource_name -f filename", api.ActionUpdate, runUpdate)
	cmd.Args = cobra.MinimumNArgs(1)
	utils.ChangeCmd(cmd, &UpdateOption.CommonParam)
	cmd.PersistentFlags().StringVarP(&UpdateOption.Filename, "filename", "f", "", "resource data file")
	cmd.MarkPersistentFlagFilename("filename", "yaml", "yml", "json")
	cmd.RegisterFlagCompletionFunc("filename", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveFilterFileExt
	})
	return cmd
}

func runUpdate(cmd *cobra.Command, args []string, spec apis.Spec) error {
	cobra.CheckErr(utils.ReadSpec(spec, UpdateOption.Filename))
	cl, err := utils.Client(log)
	if err != nil {
		return err
	}
	if UpdateOption.DryRun {
		fmt.Printf("[DryRun] update %v", spec)
		return nil
	}
	reqId, err := cl.Update(spec, nil)
	if err != nil {
		return err
	}
	if UpdateOption.Wait {
		job, err := utils.Wait(cl, reqId, UpdateOption.WaitTimeout)
		if err != nil {
			return fmt.Errorf("failed to request JobId: %s err: %w", reqId, err)
		}
		fmt.Fprintf(os.Stderr, "success request\n")
		fmt.Fprintf(os.Stderr, "JobId: %s\n", reqId)
		fmt.Fprintf(os.Stderr, "ResourceURL: %s\n", job.ResourceUrl)
	} else {
		fmt.Fprintf(os.Stderr, "success request JobId: %s\n", reqId)
	}
	return nil
}
