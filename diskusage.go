package diskusage

import (
	"fmt"
	"strings"
)

// FormatType controls the output format for disk usage.
type FormatType int

const (
	FormatRaw   FormatType = iota // bytes
	FormatHuman                   // e.g., "10.2 GB"
)

type ListOptions struct {
	FilterDriveType      DriveType
	FilterFileSystemType DriveFileSystemType
	DetectNetworkDrives  bool
}

// DriveType represents the type of drive (e.g., internal, external, network, etc.).
type DriveType string

const (
	DriveTypeInternal DriveType = "Internal"
	DriveTypeNetwork  DriveType = "Network"
)

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

// Drive represents a physical or logical disk/drive.
type Drive struct {
	Label          string // e.g., "Macintosh HD"
	Mount          string // e.g., "/"
	Total          uint64
	Used           uint64
	Free           uint64
	Percent        float64
	Type           DriveType           // e.g., Internal, External, Network, etc.
	FileSystemType DriveFileSystemType // e.g., APFS, HFS, NTFS, etc.
	Options        []string            // e.g., "rw,nosuid,nodev,noexec,relatime"
}

// HumanizeOptions allows customization of the Humanize operation.
type HumanizeOptions struct {
	WithDecimalPlaces int
	WithSuffix        bool
}

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
// If opts is provided, populates human-readable fields according to options.
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
	Get(path string) (*Drive, error)
	ListDrives() ([]Drive, error)
}

// DefaultDiskUsage is the default implementation of DiskUsage.
type DefaultDiskUsage struct{}

func (DefaultDiskUsage) Get(path string) (*Drive, error) {
	return Get(path)
}

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
