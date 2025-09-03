// Package args provides command-line argument handling utilities for the ARK SDK.
//
// This package contains functions for formatting colored output, retrieving command-line
// arguments with interactive prompts, and handling various input types including strings,
// booleans, switches, and checkbox selections.
package args

import (
	"fmt"
	"os"

	"github.com/Iilun/survey/v2"
	"github.com/cyberark/ark-sdk-golang/pkg/common"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ColorText colors the given text using the specified color.
//
// This function applies color formatting to text when coloring is enabled in the
// common package configuration. If coloring is disabled, it returns the original
// text unchanged.
//
// Parameters:
//   - text: The text string to be colored
//   - color: The color configuration to apply to the text
//
// Returns the colored text string, or the original text if coloring is disabled.
//
// Example:
//
//	red := color.New(color.FgRed)
//	coloredText := ColorText("Error message", red)
func ColorText(text string, color *color.Color) string {
	if !common.IsColoring() {
		return text
	}
	return color.Sprintf("%s", text)
}

// PrintColored prints the given text in the specified color to stdout.
//
// This function outputs colored text to stdout when interactive mode is enabled
// or output is allowed. The text will be colored only if coloring is enabled,
// otherwise it prints plain text.
//
// Parameters:
//   - text: The text content to print (can be any type that can be formatted)
//   - color: The color configuration to apply when printing
//
// The function respects the interactive and output settings from the common package
// and will not print anything if both interactive mode and output allowance are disabled.
//
// Example:
//
//	red := color.New(color.FgRed)
//	PrintColored("Error occurred", red)
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
//
// This is a convenience function that prints text using green color formatting
// to indicate successful operations or positive messages.
//
// Parameters:
//   - text: The success message to print (can be any type that can be formatted)
//
// Example:
//
//	PrintSuccess("Operation completed successfully")
func PrintSuccess(text any) {
	PrintColored(text, color.New(color.FgGreen))
}

// PrintSuccessBright prints the given text in bright green color to stdout.
//
// This is a convenience function that prints text using bold bright green color
// formatting to emphasize important successful operations or highlight positive outcomes.
//
// Parameters:
//   - text: The success message to print with emphasis (can be any type that can be formatted)
//
// Example:
//
//	PrintSuccessBright("DEPLOYMENT SUCCESSFUL")
func PrintSuccessBright(text any) {
	PrintColored(text, color.New(color.FgGreen, color.Bold))
}

// PrintFailure prints the given text in red color to stdout.
//
// This is a convenience function that prints text using red color formatting
// to indicate errors, failures, or critical issues.
//
// Parameters:
//   - text: The failure or error message to print (can be any type that can be formatted)
//
// Example:
//
//	PrintFailure("Authentication failed")
func PrintFailure(text any) {
	PrintColored(text, color.New(color.FgRed))
}

// PrintWarning prints the given text in yellow color to stdout.
//
// This is a convenience function that prints text using yellow color formatting
// to indicate warnings, cautions, or non-critical issues that require attention.
//
// Parameters:
//   - text: The warning message to print (can be any type that can be formatted)
//
// Example:
//
//	PrintWarning("Configuration file not found, using defaults")
func PrintWarning(text any) {
	PrintColored(text, color.New(color.FgYellow))
}

// PrintNormal prints the given text in default color to stdout.
//
// This is a convenience function that prints text using the default terminal color.
// It's useful for maintaining consistency with other Print functions while using
// standard formatting.
//
// Parameters:
//   - text: The message to print in normal formatting (can be any type that can be formatted)
//
// Example:
//
//	PrintNormal("Processing complete")
func PrintNormal(text any) {
	PrintColored(text, color.New())
}

// PrintNormalBright prints the given text in bright color to stdout.
//
// This is a convenience function that prints text using bold formatting with
// the default terminal color to emphasize normal messages.
//
// Parameters:
//   - text: The message to print with bold emphasis (can be any type that can be formatted)
//
// Example:
//
//	PrintNormalBright("IMPORTANT: Review the following changes")
func PrintNormalBright(text any) {
	PrintColored(text, color.New(color.Bold))
}

// GetArg retrieves the value of a command-line argument from the command flags or prompts the user for input if not provided.
//
// This function first checks if the argument is provided via command-line flags. If not found
// or if prioritizeExistingVal is true and existingVal is not empty, it uses the existing value.
// When no value is available, it prompts the user interactively for input. The function supports
// both regular text input and hidden password input.
//
// Parameters:
//   - cmd: The cobra command containing the flags to check
//   - key: The name of the flag/argument to retrieve
//   - prompt: The message to display when prompting the user for input
//   - existingVal: An existing value that may be used as default or priority
//   - hidden: Whether the input should be hidden (password-style input)
//   - prioritizeExistingVal: Whether to use existingVal over command-line flags
//   - emptyValueAllowed: Whether empty values are acceptable as valid input
//
// Returns the retrieved or inputted value and any error that occurred during the process.
// If emptyValueAllowed is false and the final value is empty, returns an error.
//
// Example:
//
//	password, err := GetArg(
//	    cmd,
//	    "password",
//	    "Enter your password:",
//	    "",
//	    true,  // hidden input
//	    false, // don't prioritize existing
//	    false, // empty not allowed
//	)
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
//
// This function first attempts to retrieve a boolean flag from the command. If prioritizeExistingVal
// is true and existingVal is not nil, it uses the existing value. When no value is determined,
// it prompts the user with a Yes/No selection interface.
//
// Parameters:
//   - cmd: The cobra command containing the boolean flags to check
//   - key: The name of the boolean flag to retrieve
//   - prompt: The message to display when prompting the user for Yes/No selection
//   - existingVal: An existing boolean value that may be used as default or priority
//   - prioritizeExistingVal: Whether to use existingVal over command-line flags
//
// Returns the boolean value (true for "Yes", false for "No") and any error that occurred
// during flag retrieval or user interaction.
//
// Example:
//
//	confirmed, err := GetBoolArg(
//	    cmd,
//	    "confirm",
//	    "Do you want to proceed?",
//	    nil,
//	    false,
//	)
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
//
// This function first checks if the switch argument is provided via command-line flags. If
// prioritizeExistingVal is true and existingVal is not empty, it uses the existing value.
// When no value is available, it prompts the user with a selection interface containing
// the provided possible values.
//
// Parameters:
//   - cmd: The cobra command containing the flags to check
//   - key: The name of the flag/argument to retrieve
//   - prompt: The message to display when prompting the user for selection
//   - possibleVals: A slice of valid options that the user can choose from
//   - existingVal: An existing value that may be used as default or priority
//   - prioritizeExistingVal: Whether to use existingVal over command-line flags
//
// Returns the selected value from the possible options and any error that occurred
// during flag retrieval or user interaction.
//
// Example:
//
//	environment, err := GetSwitchArg(
//	    cmd,
//	    "env",
//	    "Select environment:",
//	    []string{"development", "staging", "production"},
//	    "development",
//	    false,
//	)
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
//
// This function collects values from multiple command-line flags specified in the keys slice.
// If prioritizeExistingVal is true, it also includes values from the existingVals map.
// When additional selection is needed, it prompts the user with a multi-select interface
// containing the provided possible values.
//
// Parameters:
//   - cmd: The cobra command containing the flags to check
//   - keys: A slice of flag names to collect values from
//   - prompt: The message to display when prompting the user for multi-selection
//   - possibleVals: A slice of valid options that the user can choose from
//   - existingVals: A map of existing key-value pairs that may be used as defaults
//   - prioritizeExistingVal: Whether to include values from existingVals map
//
// Returns a slice of selected values and any error that occurred during flag retrieval
// or user interaction.
//
// Example:
//
//	features, err := GetCheckboxArgs(
//	    cmd,
//	    []string{"feature1", "feature2"},
//	    "Select features to enable:",
//	    []string{"auth", "logging", "metrics", "tracing"},
//	    map[string]string{"feature1": "auth"},
//	    true,
//	)
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
