package main

import (
	"github.com/Joffref/genz/cmd/genz"
	"log"
)

func main() {
	if err := genz.Execute(); err != nil {
		log.Fatal(err)
	}
}
