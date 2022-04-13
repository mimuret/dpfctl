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
		Use: "config",
	}
	cmd.AddCommand(newCmdConfigGetCurrentContext())
	cmd.AddCommand(newCmdConfigSetCurrentContext())
	cmd.AddCommand(newCmdConfigGetContext())
	cmd.AddCommand(newCmdConfigSetContext())
	return cmd
}

func newCmdConfigGetCurrentContext() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get-current-context",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(viper.GetString("current-context"))
		},
	}
	return cmd
}

func newCmdConfigSetCurrentContext() *cobra.Command {
	cmd := &cobra.Command{
		Use: "set-current-context",
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
		Use: "get-context",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := utils.GetConfig()
			if err != nil {
				return err
			}
			table := uitable.New()
			table.AddRow("name", "endpoint")
			for name, c := range cfg.Contexts {
				table.AddRow(name, c.Endpoint)
			}
			fmt.Println(table.String())
			return nil
		},
	}
	return cmd
}

func newCmdConfigSetContext() *cobra.Command {
	cmd := &cobra.Command{
		Use: "set-context name endpoint [token]",
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
