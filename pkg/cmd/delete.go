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

var DeleteOption = &deleteOpt{}

type deleteOpt struct {
	utils.CommonParam
}

func newDeleteCmd() *cobra.Command {
	cmd := utils.NewCommand("delete resource_name [id ...]", api.ActionDelete, runDelete)
	cmd.Args = cobra.MinimumNArgs(2)
	utils.ChangeCmd(cmd, &DeleteOption.CommonParam)
	return cmd
}

func runDelete(cmd *cobra.Command, args []string, spec apis.Spec) error {
	cl, err := utils.Client(log)
	if err != nil {
		return err
	}
	if DeleteOption.DryRun {
		fmt.Printf("[DryRun] delete %v", spec)
		return nil
	}
	reqId, err := cl.Delete(spec)
	if err != nil {
		return err
	}
	if DeleteOption.Wait {
		_, err := utils.Wait(cl, reqId, DeleteOption.WaitTimeout)
		if err != nil {
			return fmt.Errorf("failed to request JobId: %s err: %w", reqId, err)
		}
		log.Infof("success request\n")
		log.Infof("JobId: %s\n", reqId)
	} else {
		log.Infof("success request JobId: %s\n", reqId)
	}
	return nil
}
