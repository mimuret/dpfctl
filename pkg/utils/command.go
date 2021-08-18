package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/gosuri/uitable"
	"github.com/mimuret/dpfctl/pkg/params"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis"
	"github.com/spf13/cobra"
)

func ValidArgsFunction(action api.Action, cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return params.GetValidArgs(api.ActionRead), cobra.ShellCompDirectiveDefault
	}
	return nil, cobra.ShellCompDirectiveNoSpace
}

func CmdUsage(subcmd string, action api.Action) *uitable.Table {
	if subcmd == "get" {
		return getCmdUsage(subcmd, action)
	}
	return cmdUsage(subcmd, action)
}

func getCmdUsage(subcmd string, action api.Action) *uitable.Table {
	t := uitable.New()
	params.IterateGroup(func(groupName string, slice params.APISetSlice) {
		exist := false
		slice.IterateAPISet(func(apiSet *params.APISet) {
			var (
				readAPI *params.API
				listAPI *params.API
			)
			if apiSepc, ok := apiSet.Action[api.ActionRead]; ok {
				readAPI = &apiSepc
			}
			if apiSepc, ok := apiSet.Action[api.ActionList]; ok {
				listAPI = &apiSepc
			}
			if (readAPI != nil || listAPI != nil) && !exist {
				t.AddRow(fmt.Sprintf("\n%s API:", groupName))
				exist = true
			}
			if readAPI != nil && listAPI == nil {
				t.AddRow(fmt.Sprintf("  %s %s %s", subcmd, apiSet.Name, readAPI.Params.String()), apiSet.Description)
			} else if readAPI == nil && listAPI != nil {
				t.AddRow(fmt.Sprintf("  %s %s %s", subcmd, apiSet.Name, listAPI.Params.String()), apiSet.Description)
			} else if readAPI != nil && listAPI != nil {
				getParams := listAPI.Params
				for i := len(listAPI.Params); i < len(readAPI.Params); i++ {
					param := readAPI.Params[i]
					param.Required = false
					getParams = append(getParams, param)
				}
				t.AddRow(fmt.Sprintf("  %s %s %s", subcmd, apiSet.Name, getParams.String()), apiSet.Description)
			}

		})
	})
	return t
}

func cmdUsage(subcmd string, action api.Action) *uitable.Table {
	t := uitable.New()
	params.IterateGroup(func(groupName string, slice params.APISetSlice) {
		exist := false
		slice.IterateAPISet(func(apiSet *params.APISet) {
			if apiSepc, ok := apiSet.Action[action]; ok {
				if !exist {
					t.AddRow(fmt.Sprintf("\n%s API:", groupName))
					exist = true
				}
				t.AddRow(fmt.Sprintf("  %s %s %s", subcmd, apiSet.Name, apiSepc.Params.String()), apiSet.Description)
			}
		})
	})
	return t
}

func SetUsage(cmd *cobra.Command, action api.Action) {
	uses := strings.Split(cmd.Use, " ")
	t := CmdUsage(uses[0], action)
	cmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}
` + t.String() + `
{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)

}

func NewCommand(use string, action api.Action, runFunc func(*cobra.Command, []string, apis.Spec) error) *cobra.Command {
	cmd := &cobra.Command{
		Use: use,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return ValidArgsFunction(action, cmd, args, toComplete)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			spec, err := GetSpecFromArgs(cmd, args, action)
			if err != nil {
				return err
			}
			return runFunc(cmd, args, spec)
		},
	}
	SetUsage(cmd, action)
	return cmd
}

type CommonParam struct {
	DryRun      bool
	Wait        bool
	WaitTimeout time.Duration
}

func ChangeCmd(cmd *cobra.Command, cp *CommonParam) {
	cmd.PersistentFlags().BoolVarP(&cp.DryRun, "dry-run", "", false, "not create and update")
	cmd.PersistentFlags().BoolVarP(&cp.Wait, "wait", "", false, "wait async response")
	cmd.PersistentFlags().DurationVarP(&cp.WaitTimeout, "wait-timeout", "", time.Minute, "wait async response timeout")
}

func GetSpecFromArgs(cmd *cobra.Command, args []string, action api.Action) (apis.Spec, error) {
	cmdName := args[0]
	args = args[1:]
	apiSet := params.GetAPISetfromCmdName(cmdName)
	if apiSet == nil {
		return nil, fmt.Errorf("not support resource %s", cmdName)
	}
	apiSpec, ok := apiSet.Action[action]
	if !ok {
		return nil, fmt.Errorf("not support action %s %s", action, cmdName)
	}
	obj := apiSpec.Object.DeepCopyObject()
	spec, ok := obj.(apis.Spec)
	if !ok {
		panic(fmt.Sprintf("[BUG] failed to cast `%s` from apis.Object to api.Spec", cmdName))
	}
	if action != api.ActionCreate && action != api.ActionUpdate {
		if apiSpec.SetFunc != nil {
			if err := apiSpec.SetFunc(spec, args); err != nil {
				return nil, fmt.Errorf("failed to set arg to request: %w", err)
			}
		} else {
			if err := apiSpec.Params.SetArgs(spec, args); err != nil {
				return nil, fmt.Errorf("failed to set arg to request: %w", err)
			}
		}
	}
	return spec, nil
}
