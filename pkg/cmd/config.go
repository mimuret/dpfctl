package cmd

import (
	"fmt"

	"github.com/gosuri/uitable"
	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newCmdConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Modify dpfctl config",
	}
	cmd.AddCommand(newCmdConfigGetCurrentContext())
	cmd.AddCommand(newCmdConfigUseCurrentContext())
	cmd.AddCommand(newCmdConfigGetContext())
	cmd.AddCommand(newCmdConfigSetContext())
	cmd.AddCommand(newCmdConfigUnsetContext())
	return cmd
}

func newCmdConfigGetCurrentContext() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-context",
		Short: "Display current contexts in the dpfctl config",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(viper.GetString("current-context"))
		},
	}
	return cmd
}

func newCmdConfigUseCurrentContext() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use-context",
		Short: "Sets the current-context in the dpfctl config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := utils.GetConfig()
			if err != nil {
				return err
			}
			cfg.CurrentContext = args[0]

			if _, ok := cfg.Contexts[args[0]]; !ok {
				return fmt.Errorf("cotext `%s` not exist", args[0])
			}

			return cfg.WriteConfig()
		},
		Args: cobra.ExactArgs(1),
	}
	return cmd
}

func newCmdConfigGetContext() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-contexts",
		Short: "Display contexts in the dpfctl config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := utils.GetConfig()
			if err != nil {
				return err
			}
			table := uitable.New()
			table.AddRow("CURRENT", "NAME", "ENDPOINT")
			for name, c := range cfg.Contexts {
				use := ""
				if name == cfg.CurrentContext {
					use = "*"
				}
				table.AddRow(use, name, c.Endpoint)
			}
			fmt.Println(table.String())
			return nil
		},
	}
	return cmd
}

func newCmdConfigSetContext() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-context name endpoint [token]",
		Short: "Set context in the dpfctl config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := utils.GetConfig()
			if err != nil {
				return err
			}
			c := utils.Context{
				Endpoint: args[1],
			}
			if len(args) == 3 {
				c.Token = args[2]
			}
			cfg.Contexts[args[0]] = c

			return cfg.WriteConfig()
		},
		Args: cobra.RangeArgs(2, 3),
	}
	return cmd
}

func newCmdConfigUnsetContext() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unset-context name",
		Short: "Unsets context in the dpfctl config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := utils.GetConfig()
			if err != nil {
				return err
			}
			if cfg.CurrentContext == args[0] {
				return fmt.Errorf("%s is current context", args[0])
			}
			delete(cfg.Contexts, args[0])
			return cfg.WriteConfig()
		},
		Args: cobra.ExactArgs(1),
	}
	return cmd
}
