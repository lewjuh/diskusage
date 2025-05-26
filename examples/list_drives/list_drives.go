package examples

import (
	"fmt"
	"log"

	"github.com/lewjuh/diskusage"
)

func main() {
	drives, err := diskusage.ListDrives()
	if err != nil {
		log.Fatal(err)
	}
	for _, drive := range drives {
		total, used, free := drive.Humanize()
		fmt.Printf("Label: %s\nMount: %s\nTotal: %s\nUsed: %s\nFree: %s\nType: %s\nPercent: %f\nOptions: %s\nFileSystemType: %s\n",
			drive.Label, drive.Mount, total, used, free, drive.Type, drive.Percent, drive.Options, drive.FileSystemType)
	}
}
