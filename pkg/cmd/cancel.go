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

var CancelOption = &calcelOpt{}

type calcelOpt struct {
	utils.CommonParam
}

func newCancelCmd() *cobra.Command {
	cmd := utils.NewCommand("cancel resource_name [id ...]", api.ActionCancel, runCancel)
	cmd.Args = cobra.MinimumNArgs(2)
	utils.ChangeCmd(cmd, &CancelOption.CommonParam)
	return cmd
}

func runCancel(cmd *cobra.Command, args []string, spec apis.Spec) error {
	if CancelOption.DryRun {
		fmt.Printf("[DryRun] cancel %v", spec)
		return nil
	}
	cl, err := utils.Client(log)
	if err != nil {
		return err
	}
	reqId, err := cl.Cancel(spec)
	if err != nil {
		return err
	}
	if CancelOption.Wait {
		_, err := utils.Wait(cl, reqId, CancelOption.WaitTimeout)
		if err != nil {
			return fmt.Errorf("failed to request JobId: %s err: %w", reqId, err)
		}
		fmt.Fprintf(os.Stderr, "success request\n")
		fmt.Fprintf(os.Stderr, "JobId: %s\n", reqId)
	} else {
		fmt.Fprintf(os.Stderr, "success request JobId: %s\n", reqId)
	}
	return nil
}
