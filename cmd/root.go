/*
Copyright © 2023 Takafumi Miyanaga miya.org.0309@gmai.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tftarget",
	Short: "interactivity select resource to ( plan | appply | destroy ) with target option",
	Long: `tftarget is a CLI library for Terraform ( plan | appply | destroy ) with target option.
You can interactivity select resource to ( plan | appply | destroy ) with target option.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	s.Suffix = " loading ..."
	s.Color("green")
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
