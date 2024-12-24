/*
Copyright © 2023 Takafumi Miyanaga miya.org.0309@gmai.com
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func executePlan(cmd *cobra.Command, option string) ([]string, error) {
	action, _ := cmd.Flags().GetString("filter")
	validValues := maps.Keys(actionPatternMap)
	if err := validateFlagValue(action, validValues); action != "" && err != nil {
		return nil, err
	}

	executable, _ := cmd.Flags().GetString("executable")
	p, _ := cmd.Flags().GetInt("parallel")
	planCmd := exec.Command(executable, "plan", "-no-color", fmt.Sprintf("--parallelism=%d", p))
	if option != "" {
		planCmd = exec.Command(executable, "plan", option, "-no-color", fmt.Sprintf("--parallelism=%d", p))
	}

	out, err := planCmd.CombinedOutput()
	if err != nil {
		color.Red.Println(string(out))
		return nil, err
	}

	actionPattern := actionPatternMap[action]
	resources := extractResource(out, actionPattern)
	if len(resources) == 0 {
		return nil, ErrNotFound
	}
	options := make([]string, 0, 100)
	options = append(options, exitOpt)
	return append(options, resources...), nil
}

func extractResource(input []byte, actionPattern string) []string {
	re := regexp.MustCompile(`#\s(.*?)\s(.*)`)
	matches := re.FindAllSubmatch(input, -1)
	resourceActions := make([]string, 0, len(matches))
	if actionPattern != "" {
		for _, match := range matches {
			resource := string(match[1])
			action := string(match[2])
			if actionPattern == action {
				resourceAction := fmt.Sprintf("%s %s", resource, action)
				resourceActions = append(resourceActions, resourceAction)
			}
		}
		return resourceActions
	}
	for _, match := range matches {
		resource := string(match[1])
		action := string(match[2])
		if slices.Contains(defaultActionPattern, action) {
			resourceAction := fmt.Sprintf("%s %s", resource, action)
			resourceActions = append(resourceActions, resourceAction)
		}
	}
	return resourceActions
}

func dropAction(strs []string) []string {
	var result []string
	for _, s := range strs {
		s = strings.TrimSpace(s)
		parts := strings.Split(s, " ")
		if len(parts) > 0 {
			result = append(result, parts[0])
		}
	}
	return result
}

func slice2String(slice []string) string {
	var buffer bytes.Buffer
	if len(slice) == 1 {
		buffer.WriteString(slice[0])
		return buffer.String()
	}
	buffer.WriteString(`{'`)
	for i, item := range slice {
		buffer.WriteString(item)
		if i < len(slice)-1 {
			buffer.WriteString("','")
		}
	}
	buffer.WriteString(`'}`)
	return buffer.String()
}

func genTargetCmd(cmd *cobra.Command, action, target string) bytes.Buffer {
  var buf bytes.Buffer
  executable, _ := cmd.Flags().GetString("executable")
  buf.WriteString(executable + " " + action)
  target = strings.TrimSpace(target)
  if strings.HasPrefix(target, "{") && strings.HasSuffix(target, "}") {
    // Handle matrix of targets
    target = strings.Trim(target, "{}") // Remove surrounding braces
    targetList := strings.Split(target, ",")
    for _, t := range targetList {
      buf.WriteString(" -target=" + strings.TrimSpace(t))
      }
  } else {
    // Handle single target
    buf.WriteString(" -target=" + target)
  }
  p, _ := cmd.Flags().GetInt("parallel")
  buf.WriteString(fmt.Sprintf(" --parallelism=%d", p))
  return buf
}

func isYes(reader *bufio.Reader) bool {
	color.Red.Print("Enter a value: ")
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text) == "yes"
}

func confirm(buf bytes.Buffer) *exec.Cmd {
	buf.WriteString(" -auto-approve")
	confirm := exec.Command("sh", "-c", buf.String())
	confirm.Stdout = os.Stdout
	confirm.Stderr = os.Stderr
	return confirm
}

func targetCmd(buf bytes.Buffer) *exec.Cmd {
	cmd := exec.Command("sh", "-c", buf.String())
	cmd.Stdout = os.Stdout
	return cmd
}

func summary(items []string) error {
	color.Green.Println("==============🎉 Selected Resources 🎉==============")
	for _, v := range items {
		fmt.Println(v)
	}
	color.Green.Println("====================================================")
	return nil
}

func validateFlagValue(value string, validValues []string) error {
	for _, v := range validValues {
		if value == v {
			return nil
		}
	}
	return fmt.Errorf("invalid value: %s, valid values are: %v", value, validValues)
}
