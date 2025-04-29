package args

import (
	"fmt"
	"github.com/Iilun/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"os"
)

// ColorText colors the given text using the specified color.
func ColorText(text string, color *color.Color) string {
	if !common.IsColoring() {
		return text
	}
	return color.Sprintf("%s", text)
}

// PrintColored prints the given text in the specified color to stdout.
func PrintColored(text any, color *color.Color) {
	if common.IsInteractive() || common.IsAllowingOutput() {
		if common.IsColoring() {
			_, _ = color.Fprintf(os.Stdout, "%s\n", text)
		} else {
			_, _ = fmt.Fprintf(os.Stdout, "%s\n", text)
		}
	}
}

// PrintSuccess prints the given text in green color to stdout.
func PrintSuccess(text any) {
	PrintColored(text, color.New(color.FgGreen))
}

// PrintSuccessBright prints the given text in bright green color to stdout.
func PrintSuccessBright(text any) {
	PrintColored(text, color.New(color.FgGreen, color.Bold))
}

// PrintFailure prints the given text in red color to stdout.
func PrintFailure(text any) {
	PrintColored(text, color.New(color.FgRed))
}

// PrintWarning prints the given text in yellow color to stdout.
func PrintWarning(text any) {
	PrintColored(text, color.New(color.FgYellow))
}

// PrintNormal prints the given text in default color to stdout.
func PrintNormal(text any) {
	PrintColored(text, color.New())
}

// PrintNormalBright prints the given text in bright color to stdout.
func PrintNormalBright(text any) {
	PrintColored(text, color.New(color.Bold))
}

// GetArg retrieves the value of a command-line argument from the command flags or prompts the user for input if not provided.
func GetArg(cmd *cobra.Command, key string, prompt string, existingVal string, hidden bool, prioritizeExistingVal bool, emptyValueAllowed bool) (string, error) {
	val := ""
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Name == key {
			val = f.Value.String()
		}
	})
	if prioritizeExistingVal && existingVal != "" {
		val = existingVal
	}
	var answer string
	if hidden {
		prompt := &survey.Password{
			Message: prompt,
		}
		err := survey.AskOne(prompt, &answer)
		if err != nil {
			return "", err
		}
	} else {
		prompt := &survey.Input{
			Message: prompt,
			Default: val,
		}
		err := survey.AskOne(prompt, &answer)
		if err != nil {
			return "", err
		}
	}
	if answer != "" {
		val = answer
	}
	if val == "" && !emptyValueAllowed {
		PrintFailure("Value cannot be empty")
		return "", fmt.Errorf("value cannot be empty")
	}
	return val, nil
}

// GetBoolArg retrieves a boolean argument from the command flags or prompts the user for input if not provided.
func GetBoolArg(cmd *cobra.Command, key, prompt string, existingVal *bool, prioritizeExistingVal bool) (bool, error) {
	val := false
	if newVal, err := cmd.Flags().GetBool(key); err == nil {
		val = newVal
	}
	if prioritizeExistingVal && existingVal != nil {
		val = *existingVal
	}
	options := []string{"Yes", "No"}
	defaultOption := "No"
	if val {
		defaultOption = "Yes"
	}

	var answer string
	promptSelect := &survey.Select{
		Message: prompt,
		Options: options,
		Default: defaultOption,
	}

	err := survey.AskOne(promptSelect, &answer)
	if err != nil {
		return false, err
	}

	val = (answer == "Yes")
	return val, nil
}

// GetSwitchArg retrieves a switch argument from the command flags or prompts the user for input if not provided.
func GetSwitchArg(cmd *cobra.Command, key, prompt string, possibleVals []string, existingVal string, prioritizeExistingVal bool) (string, error) {
	val := ""
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Name == key {
			val = f.Value.String()
		}
	})
	if prioritizeExistingVal && existingVal != "" {
		val = existingVal
	}
	var answer string
	promptSelect := &survey.Select{
		Message: prompt,
		Options: possibleVals,
		Default: val,
	}

	err := survey.AskOne(promptSelect, &answer)
	if err != nil {
		return "", err
	}

	return answer, nil
}

// GetCheckboxArgs retrieves checkbox arguments from the command flags or prompts the user for input if not provided.
func GetCheckboxArgs(cmd *cobra.Command, keys []string, prompt string, possibleVals []string, existingVals map[string]string, prioritizeExistingVal bool) ([]string, error) {
	vals := []string{}
	for _, key := range keys {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Name == key {
				vals = append(vals, f.Value.String())
			}
		})
		if prioritizeExistingVal {
			if v, ok := existingVals[key]; ok {
				vals = append(vals, v)
			}
		}
	}

	selectedVals := []string{}
	promptMultiSelect := &survey.MultiSelect{
		Message: prompt,
		Options: possibleVals,
		Default: vals,
	}

	err := survey.AskOne(promptMultiSelect, &selectedVals)
	if err != nil {
		return nil, err
	}

	return selectedVals, nil
}
