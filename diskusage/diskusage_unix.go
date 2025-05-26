//go:build linux || freebsd
// +build linux freebsd

package diskusage

import (
	"syscall"
)

// getPlatformStats returns disk usage stats for Unix-like systems.
func getPlatformStats(path string) (Drive, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return Drive{}, err
	}
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used := total - free
	return Drive{Total: total, Used: used}, nil
}
