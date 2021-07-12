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

	"github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis"
	"github.com/spf13/cobra"
)

var ListOption = &listOption{}

type listOption struct {
	Output          string
	Filename        string
	RowSearchParams string
	NoHeaders       bool
}

func newListCmd() *cobra.Command {
	cmd := utils.NewCommand("list resource_name [id ...] [-o [json|yaml|line]]", api.ActionList, runList)
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.PersistentFlags().StringVarP(&ListOption.Output, "output", "o", "line", "[json|yaml|line|go-template=]")
	cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"yaml", "json", "line", "template", "go-template="}, cobra.ShellCompDirectiveFilterFileExt
	})
	cmd.PersistentFlags().StringVarP(&ListOption.Filename, "filename", "f", "", "output file name")
	cmd.PersistentFlags().StringVarP(&ListOption.Filename, "no-headers", "", "", "no output headers")

	cmd.PersistentFlags().BoolVarP(&ListOption.NoHeaders, "row-search-params", "", false, "no output header")
	return cmd
}

func runList(cmd *cobra.Command, args []string, spec apis.Spec) error {
	p, err := printer.GetPrinter(args[0], ListOption.Output)
	if err != nil {
		return fmt.Errorf("failed to get output printer: %w", err)
	}
	if linePinrer, ok := p.(*printer.LinePrinter); ok {
		linePinrer.SetNoHeaders(ListOption.NoHeaders)
	}

	cl, err := utils.Client(log)
	if err != nil {
		return err
	}
	listSpec, ok := spec.(apis.ListSpec)
	if !ok {
		return fmt.Errorf("not support list action %s", spec.GetName())
	}
	var sp api.SearchParams
	if ListOption.RowSearchParams != "" {
		sp, err = api.NewRowSearchParams(ListOption.RowSearchParams)
		if err != nil {
			return fmt.Errorf("row-search-params format error: %s", err)
		}
	}
	if countable, ok := listSpec.(api.CountableListSpec); ok {
		_, err = cl.ListALL(countable, sp)
	} else {
		_, err = cl.List(listSpec, sp)
	}
	if err != nil {
		return err
	}
	var w = os.Stdout
	if ListOption.Filename != "" {
		fp, err := os.Create(ListOption.Filename)
		if err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}
		defer fp.Close()
		w = fp
	}
	return p.Print(w, spec)
}
