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
	"fmt"
	"os"

	"github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newCmdApply() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "apply -f filename",
		Args: cobra.NoArgs,
		RunE: runApply,
	}
	utils.ChangeCmd(cmd)
	return cmd
}

func runApply(cmd *cobra.Command, args []string) error {
	var (
		results = &utils.CommandResults{}
	)
	p, err := printer.GetPrinter(results, viper.GetString("output"))
	if err != nil {
		return err
	}
	reader := utils.NewResourceReader(nil)
	resources, err := reader.GetResources(viper.GetString("filename"))
	if err != nil {
		return err
	}
	cl, err := utils.NewClient(log)
	if err != nil {
		return err
	}
	specs := []apis.Spec{}
	for _, resource := range resources {
		listSpec, ok := resource.(apis.ListSpec)
		if ok {
			for i := 0; i < listSpec.Len(); i++ {
				item := listSpec.Index(i)
				apiSpec, ok := item.(apis.Spec)
				if !ok {
					return fmt.Errorf("not support file, this spec is not apis.Spec item[%d]", i)
				}
				specs = append(specs, apiSpec)
			}
		} else {
			specs = append(specs, resource)
		}
	}
	actions := []api.Action{}
	for i, apisSpec := range specs {
		action, err := prepare(cl, apisSpec)
		if err != nil {
			return fmt.Errorf("failed to prepare item[%d]: %w", i, err)
		}
		actions = append(actions, action)
	}
	for i, apisSpec := range specs {
		if actions[i] == api.ActionUpdate {
			if viper.GetBool("dry-run") {
				fmt.Printf("[DryRun] update resource %v", apisSpec)
			} else {
				results.Add(cl.Update(context.Background(), apisSpec, nil))
			}
		}
		if actions[i] == api.ActionUpdate {
			if viper.GetBool("dry-run") {
				fmt.Printf("[DryRun] create resource %v", apisSpec)
			} else {
				results.Add(cl.Create(context.Background(), apisSpec, nil))
			}
		}
	}
	results.WaitJob(cl, viper.GetViper())
	p.Print(os.Stdout, results)
	return results.Err
}

func prepare(cl api.ClientInterface, apisSpec apis.Spec) (api.Action, error) {
	var (
		err    error
		create bool
		update bool
		exist  bool
	)
	if method, _ := apisSpec.GetPathMethod(api.ActionCreate); method != "" {
		create = true
	}
	if method, _ := apisSpec.GetPathMethod(api.ActionUpdate); method != "" {
		update = true
	}
	if !create && !update {
		return "", fmt.Errorf("this spec is not support create and update")
	}
	_, err = cl.Read(context.Background(), apisSpec)
	if api.IsNotFound(err) {
		exist = true
	}
	if exist {
		return api.ActionUpdate, nil
	}
	return api.ActionCreate, nil
}
