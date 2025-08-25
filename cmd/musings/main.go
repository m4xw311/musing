package main

import (
	"log"
	"os"

	"github.com/m4xw311/musing/cmd/musings/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}