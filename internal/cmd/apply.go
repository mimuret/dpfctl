/*
MIT License

Copyright (c) 2021 Manabu Sonoda

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

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
					return fmt.Errorf("apply file not support. This spec is not apis.Spec item[%d]", i)
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
		return "", fmt.Errorf("this spec is not supported create and updating")
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
