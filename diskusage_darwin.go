//go:build darwin
// +build darwin

package diskusage

import (
	"bytes"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
	howettplist "howett.net/plist"
)

// cString converts a NUL-terminated []byte to a Go string.
func cString(ca []byte) string {
	var buf []byte
	for _, b := range ca {
		if b == 0 {
			break
		}
		buf = append(buf, b)
	}
	return string(buf)
}

// humanFlags turns mount flags into human-readable slice.
func humanFlags(f uint64) []string {
	var out []string
	if f&unix.MNT_RDONLY != 0 {
		out = append(out, "ro")
	} else {
		out = append(out, "rw")
	}
	if f&unix.MNT_NOSUID != 0 {
		out = append(out, "nosuid")
	}
	if f&unix.MNT_NODEV != 0 {
		out = append(out, "nodev")
	}
	if f&unix.MNT_NOEXEC != 0 {
		out = append(out, "noexec")
	}
	if f&unix.MNT_DONTBROWSE != 0 {
		out = append(out, "hidden")
	}
	return out
}

// friendlyLabel generates a human-friendly name for the volume.
func friendlyLabel(mp, src, fstyp string) string {
	if mp == "/" {
		host, _ := os.Hostname()
		return host
	}
	if strings.HasPrefix(mp, "/Volumes/") {
		return filepath.Base(mp)
	}
	if u, err := url.Parse(src); err == nil && u.Scheme != "" {
		parts := strings.Split(u.Path, "/")
		return parts[len(parts)-1]
	}
	return mp
}

// getPlatformDrive returns disk usage info as a Drive for macOS systems.
func getPlatformDrive(path string) (*Drive, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, os.ErrInvalid
	}
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return nil, err
	}
	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used := total - free
	label := filepath.Base(path)
	var dtype DriveType = DriveTypeInternal

	out, err := exec.Command("diskutil", "info", "-plist", path).Output()
	if err == nil {
		var plistData map[string]interface{}
		d := howettplist.NewDecoder(bytes.NewReader(out))
		err = d.Decode(&plistData)
		if err == nil {
			if v, ok := plistData["VolumeName"].(string); ok && v != "" {
				label = v
			}
			if v, ok := plistData["Internal"].(bool); ok && v {
				dtype = DriveTypeInternal
			}
		}
	}
	return &Drive{
		Mount: path,
		Total: total,
		Used:  used,
		Free:  free,
		Label: label,
		Type:  dtype,
	}, nil
}

// ListDrives returns a list of all attached drives/disks on macOS.
func ListDrives(opts ...ListOptions) ([]Drive, error) {
	var result []Drive
	// gather stats
	n, err := unix.Getfsstat(nil, unix.MNT_NOWAIT)
	if err != nil {
		panic(err)
	}
	stats := make([]unix.Statfs_t, n)
	_, err = unix.Getfsstat(stats, unix.MNT_NOWAIT)
	if err != nil {
		panic(err)
	}

	seen := map[string]struct{}{}
	for _, fs := range stats {
		mp := cString(fs.Mntonname[:])
		src := cString(fs.Mntfromname[:])
		fstyp := cString(fs.Fstypename[:])

		// skip virtual/autofs and TimeMachine snapshots
		if fs.Flags&unix.MNT_AUTOMOUNTED != 0 {
			continue
		}
		if fstyp == "devfs" || fstyp == "autofs" {
			continue
		}
		if strings.HasPrefix(mp, "/Volumes/com.apple.TimeMachine.localsnapshots") ||
			strings.HasPrefix(mp, "/Volumes/.timemachine") {
			continue
		}
		if strings.HasPrefix(src, "/dev/disk") && strings.Contains(src, "s") && mp != "/" {
			continue
		}

		// dedupe by mount-point
		if _, ok := seen[mp]; ok {
			continue
		}
		seen[mp] = struct{}{}

		// filter network types
		isNet := fstyp == "nfs" || fstyp == "smbfs" || fstyp == "cifs" || fstyp == "webdav" || fstyp == "afpfs"
		if mp != "/" && !strings.HasPrefix(mp, "/Volumes/") && !isNet {
			continue
		}

		// calculate sizes
		bsize := uint64(fs.Bsize)
		total := fs.Blocks * bsize
		free := fs.Bavail * bsize
		used := total - (fs.Bfree * bsize)

		// prepare fields
		label := friendlyLabel(mp, src, fstyp)
		flags := humanFlags(uint64(fs.Flags))
		percent := float64(used) / float64(total) * 100

		result = append(result, Drive{
			Mount:          mp,
			Label:          label,
			Type:           DriveTypeInternal,
			FileSystemType: ParseDriveFileSystemType(fstyp),
			Total:          total,
			Used:           used,
			Free:           free,
			Percent:        percent,
			Options:        flags,
		})
	}
	return result, nil
}
