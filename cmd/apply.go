/*
Copyright © 2023 orangekame3 miya.org.0309@gmai.com
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/orangekame3/survey"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Terraform apply, interactively select resource to apply with target option",
	Long:  "Terraform apply, interactively select resource to apply with target option",
	RunE: func(cmd *cobra.Command, args []string) error {
		s.Start()
		options, err := executePlan(cmd, "")
		if err != nil && !IsNotFound(err) {
			return fmt.Errorf("plan :%w", err)
		}
		if IsNotFound(err) {
			fmt.Println(notFound)
			return nil
		}
		s.Stop()

		selected := make([]string, 0, 100)
		items, _ := cmd.Flags().GetInt("items")
		if err := survey.AskOne(&survey.MultiSelect{Message: "Select resources to target destroy:", Options: options, IgnoreCheckItem: exitOpt, HideFilter: true}, &selected, survey.WithPageSize(items)); err != nil {
			return fmt.Errorf("select resource :%w", err)
		}
		if len(selected) == 0 {
			fmt.Println(resouceNotSelected)
			return nil
		}
		if slices.Contains(selected, exitOpt) {
			fmt.Println(exitSelected)
			return nil
		}
		s.Restart()
		buf := genTargetCmd(cmd, "apply", slice2String(dropAction(selected)))
		targetCmd(buf).Run()
		s.Stop()
		if !isYes(bufio.NewReader(os.Stdin)) {
			fmt.Println(notExecuted("apply"))
			return nil
		}
		confirm(buf).Run()
		if result, _ := cmd.Flags().GetBool("summary"); result {
			summary(selected)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().IntP("parallel", "p", 10, "limit the number of concurrent operations")
	applyCmd.Flags().IntP("items", "i", 25, "check box item size")
	applyCmd.Flags().BoolP("summary", "s", true, "summary of selected items")
	applyCmd.Flags().StringP("filter", "f", "", "filter by action. You can select create, destroy, update, or replace")
}
