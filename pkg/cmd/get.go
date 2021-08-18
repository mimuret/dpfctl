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

	"github.com/mimuret/dpfctl/pkg/params"
	"github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis"
	"github.com/spf13/cobra"
)

var GetOption = &getOption{}

type getOption struct {
	Output          string
	Filename        string
	RowSearchParams string
	NoHeaders       bool
}

func newGetCmd() *cobra.Command {
	cmd := utils.NewCommand("get resource_name [id ...] [-o [json|yaml|line]]", api.ActionList, nil)
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.RunE = prepareGet
	cmd.PersistentFlags().StringVarP(&GetOption.Output, "output", "o", "line", "[json|yaml|line|go-template=]")
	cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"yaml", "json", "line", "template", "go-template="}, cobra.ShellCompDirectiveFilterFileExt
	})
	cmd.PersistentFlags().StringVarP(&GetOption.Filename, "filename", "f", "", "output file name")
	cmd.PersistentFlags().BoolVarP(&GetOption.NoHeaders, "no-headers", "", false, "no output headers")
	cmd.PersistentFlags().StringVarP(&GetOption.RowSearchParams, "row-search-params", "", "", "search params")
	return cmd
}

func prepareGet(cmd *cobra.Command, args []string) error {
	var (
		act     api.Action = api.ActionRead
		spec    apis.Spec
		readAPI *params.API
		listAPI *params.API
	)

	p, err := printer.GetPrinter(args[0], GetOption.Output)
	if err != nil {
		return fmt.Errorf("failed to get output printer: %w", err)
	}
	if linePrinter, ok := p.(*printer.LinePrinter); ok {
		linePrinter.SetNoHeaders(GetOption.NoHeaders)
	}
	var w = os.Stdout
	if GetOption.Filename != "" {
		fp, err := os.Create(GetOption.Filename)
		if err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}
		defer fp.Close()
		w = fp
	}

	resourceName := args[0]
	apiSet := params.GetAPISetfromCmdName(resourceName)
	if apiSet == nil {
		return fmt.Errorf("not support resource %s", resourceName)
	}
	paramAPI, ok := apiSet.Action[api.ActionRead]
	if ok {
		readAPI = &paramAPI
	}
	paramAPI, ok = apiSet.Action[api.ActionList]
	if ok {
		listAPI = &paramAPI
	}
	if readAPI != nil && listAPI == nil {
		act = api.ActionRead
	} else if readAPI == nil && listAPI != nil {
		act = api.ActionList
	} else if len(listAPI.Params) == len(args)-1 {
		act = api.ActionList
	}
	if spec, err = utils.GetSpecFromArgs(cmd, args, act); err != nil {
		return err
	}
	if act == api.ActionList {
		if err := runList(cmd, args, spec); err != nil {
			return err
		}
	} else {
		if err := runRead(cmd, args, spec); err != nil {
			return err
		}
	}
	return p.Print(w, spec)
}

func runList(cmd *cobra.Command, args []string, spec apis.Spec) error {
	cl, err := utils.Client(log)
	if err != nil {
		return err
	}
	listSpec, ok := spec.(apis.ListSpec)
	if !ok {
		return fmt.Errorf("not support list action %s", spec.GetName())
	}
	var sp api.SearchParams
	if GetOption.RowSearchParams != "" {
		sp, err = api.NewRowSearchParams(GetOption.RowSearchParams)
		if err != nil {
			return fmt.Errorf("row-search-params format error: %s", err)
		}
	}
	if countable, ok := listSpec.(api.CountableListSpec); ok {
		_, err = cl.ListAll(countable, sp)
	} else {
		_, err = cl.List(listSpec, sp)
	}
	if err != nil {
		return err
	}
	return nil
}

func runRead(cmd *cobra.Command, args []string, spec apis.Spec) error {
	cl, err := utils.Client(log)
	if err != nil {
		return err
	}
	_, err = cl.Read(spec)
	if err != nil {
		return err
	}
	return nil
}
