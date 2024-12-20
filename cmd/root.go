/*
Copyright © 2023 Takafumi Miyanaga miya.org.0309@gmai.com
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

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

func init() {
	executable := os.Getenv("TFTARGET_EXECUTABLE")
	if executable == "" {
		if existsInPath("terraform") {
			executable = "terraform"
		} else if existsInPath("tofu") {
			executable = "tofu"
		} else {
			fmt.Println("Error: no terraform executable found in PATH")
			os.Exit(1)
		}
	}
	rootCmd.PersistentFlags().StringP("executable", "e", executable, "The name or path of the terraform executable")
}

// existsInPath returns true if name is available in any of the paths listed
// in the PATH variable.
func existsInPath(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
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
