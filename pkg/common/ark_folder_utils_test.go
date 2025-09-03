package common

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestExpandFolder(t *testing.T) {
	// Get current user for expected results
	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("Failed to get current user: %v", err)
	}

	// Set up test environment variable
	testEnvVar := "TEST_FOLDER_VAR"
	testEnvValue := "/test/env/path"
	err = os.Setenv(testEnvVar, testEnvValue)
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	tests := []struct {
		name           string
		folder         string
		expectedResult string
	}{
		{
			name:           "empty_string",
			folder:         "",
			expectedResult: "/",
		},
		{
			name:           "simple_path_without_trailing_slash",
			folder:         "/simple/path",
			expectedResult: "/simple/path/",
		},
		{
			name:           "simple_path_with_trailing_slash",
			folder:         "/simple/path/",
			expectedResult: "/simple/path/",
		},
		{
			name:           "tilde_expansion_root_only",
			folder:         "~",
			expectedResult: currentUser.HomeDir + "/",
		},
		{
			name:           "tilde_expansion_root_with_slash",
			folder:         "~/",
			expectedResult: currentUser.HomeDir + "/",
		},
		{
			name:           "tilde_expansion_with_subpath",
			folder:         "~/Documents",
			expectedResult: filepath.Join(currentUser.HomeDir, "Documents") + "/",
		},
		{
			name:           "tilde_expansion_with_subpath_and_trailing_slash",
			folder:         "~/Documents/",
			expectedResult: filepath.Join(currentUser.HomeDir, "Documents") + "/",
		},
		{
			name:           "tilde_expansion_deep_path",
			folder:         "~/Documents/Projects/MyProject",
			expectedResult: filepath.Join(currentUser.HomeDir, "Documents/Projects/MyProject") + "/",
		},
		{
			name:           "environment_variable_expansion",
			folder:         "$" + testEnvVar,
			expectedResult: testEnvValue + "/",
		},
		{
			name:           "environment_variable_with_subpath",
			folder:         "$" + testEnvVar + "/subdir",
			expectedResult: testEnvValue + "/subdir/",
		},
		{
			name:           "environment_variable_with_braces",
			folder:         "${" + testEnvVar + "}",
			expectedResult: testEnvValue + "/",
		},
		{
			name:           "environment_variable_with_braces_and_subpath",
			folder:         "${" + testEnvVar + "}/subdir",
			expectedResult: testEnvValue + "/subdir/",
		},
		{
			name:           "non_existent_environment_variable",
			folder:         "$NON_EXISTENT_VAR",
			expectedResult: "/", // Expands to empty string, then gets trailing slash
		},
		{
			name:           "absolute_path_without_expansion",
			folder:         "/usr/local/bin",
			expectedResult: "/usr/local/bin/",
		},
		{
			name:           "relative_path",
			folder:         "relative/path",
			expectedResult: "relative/path/",
		},
		{
			name:           "path_with_spaces",
			folder:         "~/My Documents/Project Files",
			expectedResult: filepath.Join(currentUser.HomeDir, "My Documents/Project Files") + "/",
		},
		{
			name:           "path_with_special_characters",
			folder:         "~/test-folder_with.special@chars",
			expectedResult: filepath.Join(currentUser.HomeDir, "test-folder_with.special@chars") + "/",
		},
		{
			name:           "multiple_slashes",
			folder:         "//multiple//slashes//",
			expectedResult: "//multiple//slashes//",
		},
		{
			name:           "single_dot_path",
			folder:         ".",
			expectedResult: "./",
		},
		{
			name:           "double_dot_path",
			folder:         "..",
			expectedResult: "../",
		},
		{
			name:           "mixed_tilde_and_env_var_order_matters",
			folder:         "~/Documents/$" + testEnvVar,
			expectedResult: filepath.Join(currentUser.HomeDir, "Documents/"+testEnvValue) + "/",
		},
		{
			name:           "env_var_then_tilde_no_expansion",
			folder:         "$" + testEnvVar + "/~/subdir",
			expectedResult: testEnvValue + "/~/subdir/",
		},
		{
			name:           "very_long_path",
			folder:         "~/very/long/path/with/many/segments/that/goes/deep/into/filesystem/structure",
			expectedResult: filepath.Join(currentUser.HomeDir, "very/long/path/with/many/segments/that/goes/deep/into/filesystem/structure") + "/",
		},
		{
			name:           "path_with_unicode_characters",
			folder:         "~/测试/フォルダ/папка",
			expectedResult: filepath.Join(currentUser.HomeDir, "测试/フォルダ/папка") + "/",
		},
		{
			name:           "path_with_tabs_and_newlines",
			folder:         "test\tfolder\nwith\rspecial\vchars",
			expectedResult: "test\tfolder\nwith\rspecial\vchars/",
		},
		{
			name:           "path_starting_with_multiple_tildes",
			folder:         "~~double~tilde",
			expectedResult: "~~double~tilde/",
		},
		{
			name:           "tilde_in_middle_not_expanded",
			folder:         "/path/~/middle",
			expectedResult: "/path/~/middle/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ExpandFolder(tt.folder)

			if result != tt.expectedResult {
				t.Errorf("ExpandFolder(%q) = %q, want %q", tt.folder, result, tt.expectedResult)
			}
		})
	}
}
