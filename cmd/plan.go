/*
Copyright © 2023 Takafumi Miyanaga miya.org.0309@gmai.com
*/
package cmd

import (
	"fmt"

	"github.com/orangekame3/survey"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Terraform plan, interactively select resource to plan with target option",
	Long:  "Terraform plan, interactively select resource to plan with target option",
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
		if err := targetCmd(genTargetCmd(cmd, "plan", slice2String(dropAction(selected)))).Run(); err != nil {
			return err
		}
		s.Stop()
		if result, _ := cmd.Flags().GetBool("summary"); result {
			summary(selected)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.Flags().IntP("parallel", "p", 10, "limit the number of concurrent operations")
	planCmd.Flags().IntP("items", "i", 25, "check box item size")
	planCmd.Flags().BoolP("summary", "s", true, "summary of selected items")
	planCmd.Flags().StringP("filter", "f", "", "filter by action. You can select create, destroy, update, or replace")
}
