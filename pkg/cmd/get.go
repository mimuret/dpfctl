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

var GetOption = &getOption{}

type getOption struct {
	Output   string
	Filename string
}

func newGetCmd() *cobra.Command {
	cmd := utils.NewCommand("get resource_name [id ...]", api.ActionRead, runGet)
	cmd.Args = cobra.MinimumNArgs(2)
	cmd.PersistentFlags().StringVarP(&GetOption.Output, "output", "o", "line", "[json|yaml|line|go-template=]")
	cmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"yaml", "json", "line", "template", "go-template="}, cobra.ShellCompDirectiveFilterFileExt
	})
	cmd.PersistentFlags().StringVarP(&GetOption.Filename, "filename", "f", "", "output file name")
	return cmd
}

func runGet(cmd *cobra.Command, args []string, spec apis.Spec) error {
	p, err := printer.GetPrinter(args[0], GetOption.Output)
	if err != nil {
		return fmt.Errorf("failed to get output printer: %w", err)
	}

	cl, err := utils.Client(log)
	if err != nil {
		return err
	}
	_, err = cl.Read(spec)
	if err != nil {
		return err
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
	return p.Print(w, spec)
}
