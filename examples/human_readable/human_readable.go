package main

import (
	"fmt"
	"log"

	"github.com/lewjuh/diskusage"
)

// To run this example: go run human_readable.go
func main() {
	// Get drive info with human-readable formatting (2 decimal places, with suffix)
	drive, err := diskusage.Get("/")
	if err != nil {
		log.Fatal(err)
	}
	total, used, free := drive.Humanize()
	fmt.Printf("Label: %s\nMount: %s\nTotal: %s\nUsed: %s\nFree: %s\n",
		drive.Label, drive.Mount, total, used, free)
}
