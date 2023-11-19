package main

import (
	"log"

	"github.com/fornellas/tasmota_exporter/cli"
)

func main() {
	log.SetFlags(0)
	if err := cli.Cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
