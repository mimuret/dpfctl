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
	"context"

	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	cmd := utils.NewCommand("update -f filename", api.ActionUpdate, func(cmd *cobra.Command, cl api.ClientInterface, args []string, resources []apis.Spec) error {
		return commonChangeRunFunc(func(s apis.Spec) (string, error) {
			return cl.Update(context.Background(), s, nil)
		}, cmd, cl, args, resources)
	})
	cmd.Args = cobra.MinimumNArgs(0)
	utils.ChangeCmd(cmd)
	return cmd
}
