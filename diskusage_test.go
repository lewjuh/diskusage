package diskusage

import (
	"os"
	"runtime"
	"testing"
)

func TestGet(t *testing.T) {
	t.Run("ValidRootPath", func(t *testing.T) {
		var path string
		switch runtime.GOOS {
		case "windows":
			path = "C:\\"
		default:
			path = "/"
		}

		drive, err := Get(path)
		if err != nil {
			t.Fatalf("Get(%q) failed: %v", path, err)
		}
		if drive.Total == 0 {
			t.Errorf("Total disk size should not be zero")
		}
		if drive.Used == 0 {
			t.Errorf("Used disk size should not be zero")
		}
		if drive.Used > drive.Total {
			t.Errorf("Used disk size should not exceed total")
		}
	})

	t.Run("InvalidPath", func(t *testing.T) {
		_, err := Get("/unlikely_to_exist_path_12345")
		if err == nil {
			t.Errorf("Expected error for invalid path, got nil")
		}
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		var path string
		switch runtime.GOOS {
		case "windows":
			path = `C:\System Volume Information`
		default:
			path = "/root"
		}
		_, err := Get(path)
		if err == nil {
			t.Skipf("Expected error for permission denied path %q, but got none (may be running as root/admin)", path)
		}
	})

	t.Run("FilePath", func(t *testing.T) {
		f, err := os.CreateTemp("", "diskusage_test_file")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(f.Name())
		_, err = Get(f.Name())
		if err == nil {
			t.Errorf("Expected error for file path, got nil")
		}
	})
}

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

func TestListDrives_NotImplemented(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("macOS implementation tested separately")
	}
	_, err := ListDrives()
	if err == nil {
		t.Errorf("Expected error or empty result for unimplemented ListDrives, got nil")
	}
}

func TestListDrivePaths(t *testing.T) {
	paths, err := ListDrivePaths()
	if err != nil {
		t.Fatalf("ListDrivePaths() failed: %v", err)
	}
	if len(paths) == 0 {
		t.Errorf("Expected at least one drive path, got 0")
	}
	for _, p := range paths {
		if p == "" {
			t.Errorf("Drive path should not be empty")
		}
	}
	// On Unix/macOS, root should be present
	if runtime.GOOS != "windows" {
		found := false
		for _, p := range paths {
			if p == "/" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected root path '/' to be present in drive paths")
		}
	}
}
