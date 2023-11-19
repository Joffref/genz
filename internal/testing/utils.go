package testing

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"os"
)

// assertOutputIsEqual compares the expected and actual byte slices and returns an error if they are different.
// If verbose is true, the output is printed. This is useful for debugging.
func assertOutputIsEqual(expected, actual []byte, verbose bool) error {
	if string(expected) != string(actual) {
		if verbose {
			return fmt.Errorf("Difference between expected.go and generated file:\n%s\n", cmp.Diff(expected, actual))
		}
		return fmt.Errorf("expected.go and generated file are different")
	}
	return nil
}

func removeFiles(generatedFiles []string) error {
	for _, f := range generatedFiles {
		if err := os.Remove(f); err != nil {
			return fmt.Errorf("failed to remove file %s: %w", f, err)
		}
	}
	return nil
}
