// Package diskusage provides a simple, cross-platform API for retrieving disk usage statistics.
// It supports querying total, used, and free space for a given path, as well as listing all drives on supported platforms.
package diskusage

import (
	"fmt"
	"strings"
)

// FormatType controls the output format for disk usage.
type FormatType int

const (
	// FormatRaw outputs values in bytes.
	FormatRaw FormatType = iota // bytes
	// FormatHuman outputs values in a human-readable format (e.g., "10.2 GB").
	FormatHuman // e.g., "10.2 GB"
)

// ListOptions specifies options for listing drives.
type ListOptions struct {
	FilterDriveType      DriveType
	FilterFileSystemType DriveFileSystemType
	DetectNetworkDrives  bool
}

// DriveType represents the type of drive (e.g., internal, external, network, etc.).
type DriveType string

const (
	// DriveTypeInternal represents an internal drive.
	DriveTypeInternal DriveType = "Internal"
	// DriveTypeNetwork represents a network drive.
	DriveTypeNetwork DriveType = "Network"
)

// DriveFileSystemType represents the file system type of a drive.
type DriveFileSystemType string

const (
	DriveFileSystemTypeAPFS    DriveFileSystemType = "APFS"
	DriveFileSystemTypeHFS     DriveFileSystemType = "HFS"
	DriveFileSystemTypeExFAT   DriveFileSystemType = "ExFAT"
	DriveFileSystemTypeNTFS    DriveFileSystemType = "NTFS"
	DriveFileSystemTypeFAT32   DriveFileSystemType = "FAT32"
	DriveFileSystemTypeExt4    DriveFileSystemType = "Ext4"
	DriveFileSystemTypeNFS     DriveFileSystemType = "NFS"
	DriveFileSystemTypeSMBFS   DriveFileSystemType = "SMBFS"
	DriveFileSystemTypeAFPFS   DriveFileSystemType = "AFPFS"
	DriveFileSystemTypeAutofs  DriveFileSystemType = "Autofs"
	DriveFileSystemTypeWebDAV  DriveFileSystemType = "WebDAV"
	DriveFileSystemTypeSSHFS   DriveFileSystemType = "SSHFS"
	DriveFileSystemTypeUnknown DriveFileSystemType = "Unknown"
)

// Drive represents a physical or logical disk/drive and its usage statistics.
type Drive struct {
	Label          string              // e.g., "Macintosh HD"
	Mount          string              // e.g., "/"
	Total          uint64              // Total bytes
	Used           uint64              // Used bytes
	Free           uint64              // Free bytes
	Percent        float64             // Used percentage
	Type           DriveType           // e.g., Internal, External, Network, etc.
	FileSystemType DriveFileSystemType // e.g., APFS, HFS, NTFS, etc.
	Options        []string            // e.g., "rw,nosuid,nodev,noexec,relatime"
}

// HumanizeOptions allows customization of the Humanize operation.
type HumanizeOptions struct {
	WithDecimalPlaces int  // Number of decimal places to use
	WithSuffix        bool // Whether to include a unit suffix (e.g., GB)
}

// Humanize returns human-readable strings for total, used, and free space.
// You can optionally provide HumanizeOptions to control formatting.
func (d *Drive) Humanize(opts ...HumanizeOptions) (string, string, string) {
	decimals := 2
	suffix := true
	if len(opts) > 0 {
		decimals = opts[0].WithDecimalPlaces
		if decimals == 0 {
			decimals = 1
		}
		suffix = opts[0].WithSuffix
	}
	totalHumanized := humanizeBytes(d.Total, decimals, suffix)
	usedHumanized := humanizeBytes(d.Used, decimals, suffix)
	freeHumanized := humanizeBytes(d.Free, decimals, suffix)
	return totalHumanized, usedHumanized, freeHumanized
}

// Get returns disk usage information for the given path as a Drive.
// If the path is invalid or not a directory, an error is returned.
func Get(path string) (*Drive, error) {
	d, err := getPlatformDrive(path)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// humanizeBytes returns a human-readable string for a byte count.
func humanizeBytes(b uint64, decimals int, suffix bool) string {
	const unit = 1024
	if b < unit {
		if suffix {
			return fmt.Sprintf("%d B", b)
		}
		return fmt.Sprintf("%d", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	format := fmt.Sprintf("%%.%df", decimals)
	val := fmt.Sprintf(format, float64(b)/float64(div))
	if suffix {
		return fmt.Sprintf("%s %cB", val, "KMGTPEZY"[exp])
	}
	return val
}

// DiskUsage defines the interface for disk usage operations.
type DiskUsage interface {
	// Get returns disk usage information for the given path as a Drive.
	Get(path string) (*Drive, error)
	// ListDrives returns a slice of all attached drives.
	ListDrives() ([]Drive, error)
}

// DefaultDiskUsage is the default implementation of DiskUsage.
type DefaultDiskUsage struct{}

// Get returns disk usage information for the given path as a Drive (DefaultDiskUsage implementation).
func (DefaultDiskUsage) Get(path string) (*Drive, error) {
	return Get(path)
}

// ListDrives returns a slice of all attached drives (DefaultDiskUsage implementation).
func (DefaultDiskUsage) ListDrives() ([]Drive, error) {
	return ListDrives()
}

// New returns the default DiskUsage implementation.
func New() DiskUsage {
	return DefaultDiskUsage{}
}

var driveFileSystemTypeMap = map[string]DriveFileSystemType{
	"apfs":   DriveFileSystemTypeAPFS,
	"hfs":    DriveFileSystemTypeHFS,
	"exfat":  DriveFileSystemTypeExFAT,
	"ntfs":   DriveFileSystemTypeNTFS,
	"fat32":  DriveFileSystemTypeFAT32,
	"ext4":   DriveFileSystemTypeExt4,
	"nfs":    DriveFileSystemTypeNFS,
	"smbfs":  DriveFileSystemTypeSMBFS,
	"afpfs":  DriveFileSystemTypeAFPFS,
	"autofs": DriveFileSystemTypeAutofs,
	"webdav": DriveFileSystemTypeWebDAV,
	"sshfs":  DriveFileSystemTypeSSHFS,
}

// ParseDriveFileSystemType returns the DriveFileSystemType constant matching s (case-insensitive), or DriveFileSystemTypeUnknown if not found.
func ParseDriveFileSystemType(s string) DriveFileSystemType {
	if t, ok := driveFileSystemTypeMap[strings.ToLower(s)]; ok {
		return t
	}
	return DriveFileSystemTypeUnknown
}
