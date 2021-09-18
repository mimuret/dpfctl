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
	"fmt"
	"os"

	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var log = &api.StdLogger{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dpfctl",
	Short: "dpfctl controls IIJ DNS Platform Service",
	Long:  `A`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if _, err = utils.GetContexts(); err != nil {
			return err
		}
		log.LogLevel = viper.GetInt("debug")
		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.dpfctl.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.PersistentFlags().Int("debug", 2, "debug level 0=trace,1=debug,2=info,3=error (default is 2)")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().String("context", "", "current context")
	viper.BindPFlag("context", rootCmd.PersistentFlags().Lookup("context"))

	rootCmd.PersistentFlags().StringP("output", "o", "line", "[json|yaml|line|go-template=]")
	rootCmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"yaml", "json", "line", "template", "go-template="}, cobra.ShellCompDirectiveFilterFileExt
	})
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))

	rootCmd.PersistentFlags().BoolP("no-headers", "", false, "no output headers")
	viper.BindPFlag("no-headers", rootCmd.PersistentFlags().Lookup("no-headers"))

	rootCmd.AddCommand(newGetCmd())

	rootCmd.AddCommand(newCreateCmd())
	rootCmd.AddCommand(newUpdateCmd())
	rootCmd.AddCommand(newDeleteCmd())
	rootCmd.AddCommand(newCancelCmd())
	rootCmd.AddCommand(newRunCmd())

	rootCmd.AddCommand(newCmdApply())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if viper.GetString("config") != "" {
		// Use config file from the flag.
		viper.SetConfigFile(viper.GetString("config"))
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cmd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dpfctl")
	}
	viper.SetEnvPrefix("dpf")
	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}