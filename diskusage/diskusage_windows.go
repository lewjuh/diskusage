//go:build windows
// +build windows

package diskusage

import (
	"syscall"
	"unsafe"
)

// getPlatformStats returns disk usage stats for Windows systems.
func getPlatformStats(path string) (Drive, error) {
	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)

	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return Drive{}, err
	}

	r, _, err := syscall.NewLazyDLL("kernel32.dll").NewProc("GetDiskFreeSpaceExW").Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)),
	)
	if r == 0 {
		return Drive{}, err
	}

	used := lpTotalNumberOfBytes - lpTotalNumberOfFreeBytes
	return Drive{Total: uint64(lpTotalNumberOfBytes), Used: uint64(used)}, nil
}
