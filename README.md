# diskusage

A simple, cross-platform Go package for retrieving disk usage statistics (total and used space).

## Installation

```
go get github.com/lewjuh/diskusage
```

## Supported Platforms

- **macOS (Darwin)**: Full support (Get, ListDrives)
- **Linux/Unix (Linux, FreeBSD, etc.)**: Full support (Get, ListDrives)
- **Windows**: Only `Get` is supported; `ListDrives` is not implemented

## Compatibility

| Feature      | macOS | Linux/Unix | Windows |
|--------------|:-----:|:----------:|:-------:|
| `Get`        |  ✔️   |    ✔️      |   ❌    |
| `ListDrives` |  ✔️   |    ✔️      |   ❌    |

Windows compatibility coming soon.

## API

### Functions
- `Get(path string) (*Drive, error)`
  - Returns disk usage information for the given path as a `*Drive`.
- `ListDrives(opts ...ListOptions) ([]Drive, error)`
  - Returns a slice of all attached drives (not implemented on Windows).

### Types
- `type Drive struct`:
  - `Label string`           // e.g., "Macintosh HD"
  - `Mount string`           // e.g., "/"
  - `Total uint64`           // Total bytes
  - `Used uint64`            // Used bytes
  - `Free uint64`            // Free bytes
  - `Percent float64`        // Used percentage
  - `Type DriveType`         // Internal, Network, etc.
  - `FileSystemType DriveFileSystemType` // APFS, NTFS, etc.
  - `Options []string`       // Mount options
- `func (d *Drive) Humanize(opts ...HumanizeOptions) (total, used, free string)`
  - Returns human-readable strings for total, used, and free space.

## Usage

### Get disk usage for a path
```go
import "github.com/lewjuh/diskusage"

drive, err := diskusage.Get("/")
if err != nil {
    // handle error
}
total, used, free := drive.Humanize()
fmt.Printf("Label: %s\nMount: %s\nTotal: %s\nUsed: %s\nFree: %s\nType: %s\n", drive.Label, drive.Mount, total, used, free, drive.Type)
```

### List all drives 
```go
import "github.com/lewjuh/diskusage"

drives, err := diskusage.ListDrives()
if err != nil {
    // handle error
}
for _, drive := range drives {
    total, used, free := drive.Humanize()
    fmt.Printf("Label: %s\nMount: %s\nTotal: %s\nUsed: %s\nFree: %s\nType: %s\nPercent: %.2f%%\nOptions: %v\nFileSystemType: %s\n\n",
        drive.Label, drive.Mount, total, used, free, drive.Type, drive.Percent, drive.Options, drive.FileSystemType)
}
``` 