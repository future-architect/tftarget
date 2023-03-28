/*
Copyright © 2023 Takafumi Miyanaga miya.org.0309@gmai.com
*/
package cmd

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/gookit/color"
)

var (
	resouceNotSelected   = color.Green.Sprint("resource not seleced")
	exitSelected         = color.Green.Sprintf("exit seleced")
	notFound             = color.Green.Sprintf("not found target resource")
	s                    = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	exitOpt              = color.Red.Sprintf("%s", "exit (cancel terraform plan)")
	defaultActionPattern = []string{"will be created", "will be destroyed", "will be updated in-place", "must be replaced"}
	actionPatternMap     = map[string]string{
		"create":  "will be created",
		"destroy": "will be destroyed",
		"update":  "will be updated in-place",
		"replace": "must be replaced",
	}
)

func notExecuted(action string) string {
	return color.Green.Sprintf("%s exit seleced", action)
}
