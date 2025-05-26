# diskusage

A simple, cross-platform Go package for retrieving disk usage statistics (total and used space).

## Usage

```go
import "github.com/lewjuh/diskusage"

stats, err := diskusage.Get("/")
if err != nil {
    // handle error
}
fmt.Printf("Total: %d, Used: %d\n", stats.Total, stats.Used)
``` 