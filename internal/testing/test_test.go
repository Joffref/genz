package testing

import (
	"runtime"
	"testing"
)

func TestRunTests(t *testing.T) {
	if runtime.GOOS == "windows" { // TODO: Fix this test on windows
		t.Skip("Skipping test on windows")
	}
	tests := map[string]struct {
		directory string
		wantErr   bool
	}{
		"Non existing directory should return an error": {
			directory: "",
			wantErr:   true,
		},
		"Directory with no expected.go file should not return an error": {
			directory: "./testdata/empty_directory",
			wantErr:   false,
		},
		"Directory with expected.go file should not return an error": {
			directory: "./testdata/accurate_expected",
			wantErr:   false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if err := RunTests(tc.directory, false, false); (err != nil) != tc.wantErr {
				t.Errorf("RunTests() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func Test_runTest(t *testing.T) {
	if runtime.GOOS == "windows" { // TODO: fix this test on windows
		t.Skip("Skipping test on windows")
	}
	tests := map[string]struct {
		directory string
		wantErr   bool
	}{
		"Non existing directory should return an error": {
			directory: "",
			wantErr:   true,
		},
		"Directory with no expected.go file should return an error": {
			directory: "./testdata/empty_directory",
			wantErr:   true,
		},
		"Directory with expected.go and identical generated file should not return an error": {
			directory: "./testdata/accurate_expected",
			wantErr:   false,
		},
		"Directory with two generated files should return an error": {
			directory: "./testdata/two_generated_files",
			wantErr:   true,
		},
		"Directory with difference between expected.go and generated file should return an error": {
			directory: "./testdata/inaccurate_expected",
			wantErr:   true,
		},
		"Directory with failing test should return an error": {
			directory: "./testdata/failing_unit_test",
			wantErr:   true,
		},
		"Directory with working test should not return an error": {
			directory: "./testdata/passing_unit_test",
			wantErr:   false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if err := runTest(tc.directory, false); (err != nil) != tc.wantErr {
				t.Errorf("runTest() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
