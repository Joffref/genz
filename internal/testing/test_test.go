package testing

import "testing"

func TestRunTests(t *testing.T) {
	type args struct {
		directory   string
		verbose     bool
		exitOnError bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Non existing directory should return an error",
			args:    args{directory: "", verbose: false, exitOnError: false},
			wantErr: true,
		},
		{
			name:    "Directory with no expected.go file should not return an error",
			args:    args{directory: "testdata/empty_directory", verbose: false, exitOnError: false},
			wantErr: false,
		},
		{
			name:    "Directory with expected.go file should not return an error",
			args:    args{directory: "./testdata/accurate_expected", verbose: false, exitOnError: false},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunTests(tt.args.directory, tt.args.verbose, tt.args.exitOnError); (err != nil) != tt.wantErr {
				t.Errorf("RunTests() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_runTest(t *testing.T) {
	type args struct {
		directory string
		verbose   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Non existing directory should return an error",
			args:    args{directory: "", verbose: false},
			wantErr: true,
		},
		{
			name:    "Directory with no expected.go file should return an error",
			args:    args{directory: "./testdata/empty_directory", verbose: false},
			wantErr: true,
		},
		{
			name:    "Directory with expected.go and identical generated file should not return an error",
			args:    args{directory: "./testdata/accurate_expected", verbose: false},
			wantErr: false,
		},
		{
			name:    "Directory with two generated files should return an error",
			args:    args{directory: "./testdata/two_generated_files", verbose: false},
			wantErr: true,
		},
		{
			name:    "Directory with difference between expected.go and generated file should return an error",
			args:    args{directory: "./testdata/difference_with_expected", verbose: false},
			wantErr: true,
		},
		{
			name:    "Directory with failing test should return an error",
			args:    args{directory: "./testdata/failing_unit_test", verbose: false},
			wantErr: true,
		},
		{
			name:    "Directory with working test should not return an error",
			args:    args{directory: "./testdata/passing_unit_test", verbose: false},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := runTest(tt.args.directory, tt.args.verbose); (err != nil) != tt.wantErr {
				t.Errorf("runTest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
