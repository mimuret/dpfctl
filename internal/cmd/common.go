package cmd

import (
	"fmt"
	"os"

	"github.com/mimuret/dpfctl/pkg/printer"
	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func commonChangeRunFunc(runFunc func(apis.Spec) (string, error), cmd *cobra.Command, cl api.ClientInterface, args []string, resources []apis.Spec) error {
	p, err := printer.GetPrinter(args[0], viper.GetString("output"))
	if err != nil {
		return err
	}
	results := &utils.CommandResults{}
	for _, resource := range resources {
		if viper.GetBool("dry-run") {
			fmt.Printf("[DryRun] %s %v", args[0], resource)
		} else {
			results.Add(runFunc(resource))
		}
	}
	if !viper.GetBool("dry-run") {
		results.WaitJob(cl, viper.GetViper())
	}
	p.Print(os.Stdout, results)
	return results.Err
}
