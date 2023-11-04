package utils

import (
	"log"
	"os"
)

// IsDirectory reports whether the named file is a directory.
func IsDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
