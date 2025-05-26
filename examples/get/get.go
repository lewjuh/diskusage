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
	total, used, free := drive.Humanize()
	fmt.Printf("Label: %s\nMount: %s\nTotal: %s\nUsed: %s\nFree: %s\nType: %s\n",
		drive.Label, drive.Mount, total, used, free, drive.Type)
}
