package diskusage

import (
	"runtime"
	"testing"
)

func TestListDrives_macOS(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("skipping macOS-specific test")
	}

	drives, err := ListDrives()
	if err != nil {
		t.Fatalf("ListDrives() failed: %v", err)
	}
	if len(drives) == 0 {
		t.Errorf("Expected at least one drive, got 0")
	}
	mounts := make(map[string]struct{})
	labels := make(map[string]struct{})
	for _, d := range drives {
		if d.Mount == "" {
			t.Errorf("Drive mount point should not be empty")
		}
		if d.Total == 0 {
			t.Errorf("Drive total size should not be zero")
		}
		if d.Free == 0 {
			t.Errorf("Drive free size should not be zero")
		}
		if d.Used > d.Total {
			t.Errorf("Drive used size should not exceed total")
		}
		if d.Label == "" {
			t.Errorf("Drive label should not be empty")
		}
		if _, ok := mounts[d.Mount]; ok {
			t.Errorf("Duplicate mount point: %q", d.Mount)
		}
		mounts[d.Mount] = struct{}{}
		if _, ok := labels[d.Label]; ok {
			// Labels can be duplicated, but warn
			t.Logf("Warning: duplicate label: %q", d.Label)
		}
		labels[d.Label] = struct{}{}
		if d.Used+d.Free < d.Total-1024*1024 && d.Used+d.Free > d.Total+1024*1024 {
			t.Errorf("Used + Free should be approximately Total (got Used=%d, Free=%d, Total=%d)", d.Used, d.Free, d.Total)
		}
		// New checks for Type field
		if d.Type == "" {
			t.Errorf("Drive type should not be empty for mount %q", d.Mount)
		}
		switch d.Type {
		case DriveTypeInternal, DriveTypeNetwork:
			// valid
		default:
			t.Errorf("Drive type %q is not a valid DriveType for mount %q", d.Type, d.Mount)
		}
		if d.Options == nil {
			t.Errorf("Options field should not be nil for mount %q", d.Mount)
		}
		if runtime.GOOS == "darwin" && len(d.Options) == 0 {
			t.Errorf("Options field should not be empty for mount %q on darwin", d.Mount)
		}
	}
}

func TestListDrives_Unix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping Unix-specific test")
	}

	drives, err := ListDrives()
	if err != nil {
		t.Fatalf("ListDrives() failed: %v", err)
	}
	if len(drives) == 0 {
		t.Errorf("Expected at least one drive, got 0")
	}
	mounts := make(map[string]struct{})
	labels := make(map[string]struct{})
	for _, d := range drives {
		if d.Mount == "" {
			t.Errorf("Drive mount point should not be empty")
		}
		if d.Total == 0 {
			t.Errorf("Drive total size should not be zero")
		}
		if d.Free == 0 {
			t.Errorf("Drive free size should not be zero")
		}
		if d.Used > d.Total {
			t.Errorf("Drive used size should not exceed total")
		}
		if d.Label == "" {
			t.Errorf("Drive label should not be empty")
		}
		if _, ok := mounts[d.Mount]; ok {
			t.Errorf("Duplicate mount point: %q", d.Mount)
		}
		mounts[d.Mount] = struct{}{}
		if _, ok := labels[d.Label]; ok {
			// Labels can be duplicated, but warn
			t.Logf("Warning: duplicate label: %q", d.Label)
		}
		labels[d.Label] = struct{}{}
		if d.Used+d.Free < d.Total-1024*1024 && d.Used+d.Free > d.Total+1024*1024 {
			t.Errorf("Used + Free should be approximately Total (got Used=%d, Free=%d, Total=%d)", d.Used, d.Free, d.Total)
		}
		// New checks for Type field
		if d.Type == "" {
			t.Errorf("Drive type should not be empty for mount %q", d.Mount)
		}
		switch d.Type {
		case DriveTypeInternal, DriveTypeNetwork:
			// valid
		default:
			t.Errorf("Drive type %q is not a valid DriveType for mount %q", d.Type, d.Mount)
		}
		if d.Options == nil {
			t.Errorf("Options field should not be nil for mount %q", d.Mount)
		}
		if runtime.GOOS == "darwin" && len(d.Options) == 0 {
			t.Errorf("Options field should not be empty for mount %q on darwin", d.Mount)
		}
	}
}
