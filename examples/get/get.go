package main

import (
	"fmt"
	"log"

	"github.com/lewjuh/diskusage"
)

func main() {
	// Get drive information for a specific path
	drive, err := diskusage.Get("/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Label: %s\nMount: %s\nTotal: %d\nUsed: %d\nFree: %d\nType: %s\n",
		drive.Label, drive.Mount, drive.Total, drive.Used, drive.Free, drive.Type)
}
