//go:build linux || freebsd
// +build linux freebsd

package diskusage

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

// skipFSTypes lists virtual or unwanted filesystems to ignore
var skipFSTypes = map[string]struct{}{
	"proc": {}, "sysfs": {}, "tmpfs": {}, "devtmpfs": {}, "devpts": {},
	"cgroup": {}, "cgroup2": {}, "securityfs": {}, "pstore": {},
	"efivarfs": {}, "debugfs": {}, "tracefs": {}, "rpc_pipefs": {},
	"overlay": {}, "squashfs": {},
}

// mountEntry holds parsed /proc/mounts data
type mountEntry struct {
	src, mp, fstyp, opts string
}

var ErrNotImplemented = errors.New("not implemented on this platform")

// getPlatformStats returns disk usage stats for Unix-like systems.
func getPlatformDrive(path string) (*Drive, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, os.ErrInvalid
	}

	drives, err := ListDrives()
	if err != nil {
		return nil, err
	}
	for _, d := range drives {
		if d.Mount == path {
			return &d, nil
		}
	}
	return nil, os.ErrInvalid
}

func ListDrives(opts ...ListOptions) ([]Drive, error) {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	seenSrc := make(map[string]struct{})
	var drives []Drive
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}
		src, mp, fstyp, opts := fields[0], fields[1], fields[2], fields[3]
		if _, ok := seenSrc[src]; ok {
			continue
		}
		seenSrc[src] = struct{}{}
		if _, skip := skipFSTypes[fstyp]; skip {
			continue
		}

		var fs unix.Statfs_t
		if err := unix.Statfs(mp, &fs); err != nil {
			continue
		}
		bsize := uint64(fs.Frsize)
		total := fs.Blocks * bsize
		free := fs.Bavail * bsize
		used := total - free
		if total == 0 && used == 0 && free == 0 {
			continue
		}

		fsType := ParseDriveFileSystemType(fstyp)
		isNet := fstyp == "nfs" || fstyp == "smbfs" || fstyp == "cifs" || fstyp == "webdav" || fstyp == "afpfs"

		dType := DriveTypeInternal
		if isNet {
			dType = DriveTypeNetwork
		}

		label := filepath.Base(src)
		if label == "" || label == "/" {
			label = filepath.Base(mp)
		}

		drive := Drive{
			Label:          label,
			Mount:          mp,
			Total:          total,
			Used:           used,
			Free:           free,
			Percent:        100 * float64(used) / float64(total),
			Type:           dType,
			FileSystemType: fsType,
			Options:        strings.Split(opts, ","),
		}
		drives = append(drives, drive)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return drives, nil
}
