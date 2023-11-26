package testing

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"log"
	"os"
)

// assertOutputIsEqual compares the expected and actual byte slices and returns an error if they are different.
// If verbose is true, the output is printed. This is useful for debugging.
func assertOutputIsEqual(expectedFileName, actualFileName string, expected, actual []byte, verbose bool) error {
	if verbose {
		log.Printf("Comparing %s and %s\n", expectedFileName, actualFileName)
	}
	if string(expected) != string(actual) {
		if verbose {
			return fmt.Errorf("Difference between %s and %s:\n%s", expectedFileName, actualFileName, cmp.Diff(string(expected), string(actual)))
		}
		return fmt.Errorf("Difference between %s and %s:\n", expectedFileName, actualFileName)
	}
	if verbose {
		log.Printf("No difference between %s and %s\n", expectedFileName, actualFileName)
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
