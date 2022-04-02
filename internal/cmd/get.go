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
	cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return utils.ValidArgsFunction(api.ActionList, args)
	}
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
	p, err := printer.GetPrinter(spec, GetOption.Output)
	if err != nil {
		return fmt.Errorf("failed to get output printer: %w", err)
	}
	if linePrinter, ok := p.(*printer.LinePrinter); ok {
		linePrinter.SetNoHeader(GetOption.NoHeaders)
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
	cl, err := utils.NewClient(log)
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
		_, err = cl.ListAll(context.Background(), countable, sp)
	} else {
		_, err = cl.List(context.Background(), listSpec, sp)
	}
	if err != nil {
		return err
	}
	return nil
}

func runRead(cmd *cobra.Command, args []string, spec apis.Spec) error {
	cl, err := utils.NewClient(log)
	if err != nil {
		return err
	}
	_, err = cl.Read(context.Background(), spec)
	if err != nil {
		return err
	}
	return nil
}
