//go:build !darwin
// +build !darwin

package diskusage

import "errors"

var ErrNotImplemented = errors.New("not implemented on this platform")

func ListDrives() ([]Drive, error) {
	return nil, ErrNotImplemented
}

func getPlatformDrive(path string) (Drive, error) {
	return Drive{}, ErrNotImplemented
}

func ListDrivePaths() ([]string, error) {
	return nil, ErrNotImplemented
}
