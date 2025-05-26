package examples

import (
	"fmt"
	"log"

	"github.com/lewjuh/diskusage"
)

func main() {
	paths, err := diskusage.ListDrivePaths()
	if err != nil {
		log.Fatal(err)
	}
	for _, path := range paths {
		fmt.Println(path)
	}
}
