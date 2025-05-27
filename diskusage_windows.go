//go:build windows
// +build windows

package diskusage

import (
	"errors"
)

var ErrNotImplemented = errors.New("not implemented on this platform")

// getPlatformDrive returns disk usage stats for Windows systems.
func getPlatformDrive(path string) (*Drive, error) {
	return nil, ErrNotImplemented
}

func ListDrives(opts ...ListOptions) ([]Drive, error) {
	return nil, ErrNotImplemented
}

func ListDrivePaths() ([]string, error) {
	return nil, ErrNotImplemented
}
