package fs

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestParentDirWithRelativesUnix(t *testing.T) {
	t.Skip("these do not work yet")
	// Helper function to create platform-specific paths
	toOSPath := func(path string) string {
		if runtime.GOOS == "windows" {
			// Convert forward slashes to OS-specific separator
			return filepath.FromSlash(path)
		}
		return path
	}

	// Helper function to create platform-specific PathInfo
	expectedPathInfo := func(parent string, relatives []string) PathInfo {
		osParent := toOSPath(parent)
		osRelatives := make([]string, len(relatives))
		for i, rel := range relatives {
			osRelatives[i] = toOSPath(rel)
		}
		return PathInfo{
			CommonParent:  osParent,
			RelativePaths: osRelatives,
		}
	}

	tests := []struct {
		name     string
		paths    []string
		expected PathInfo
		wantErr  bool
	}{
		{
			name:     "empty input",
			paths:    []string{},
			expected: PathInfo{},
			wantErr:  false,
		},
		{
			name:  "single file",
			paths: []string{"/home/user/file.txt"},
			expected: expectedPathInfo(
				"/home/user",
				[]string{"file.txt"},
			),
			wantErr: false,
		},
		{
			name: "simple case - same directory",
			paths: []string{
				"/home/user/docs/file1.txt",
				"/home/user/docs/file2.txt",
			},
			expected: expectedPathInfo(
				"/home/user/docs",
				[]string{"file1.txt", "file2.txt"},
			),
			wantErr: false,
		},
		{
			name: "nested directories",
			paths: []string{
				"/home/user/docs/file1.txt",
				"/home/user/docs/subdir/file2.txt",
				"/home/user/docs/subdir/deeper/file3.txt",
			},
			expected: expectedPathInfo(
				"/home/user/docs",
				[]string{
					"file1.txt",
					"subdir/file2.txt",
					"subdir/deeper/file3.txt",
				},
			),
			wantErr: false,
		},
		{
			name: "different depths",
			paths: []string{
				"/home/user/docs/deep/deeper/deepest/file1.txt",
				"/home/user/docs/file2.txt",
			},
			expected: expectedPathInfo(
				"/home/user/docs",
				[]string{
					"deep/deeper/deepest/file1.txt",
					"file2.txt",
				},
			),
			wantErr: false,
		},
		{
			name: "no common parent except root",
			paths: []string{
				"/home/user1/file1.txt",
				"/home/user2/file2.txt",
				"/var/log/file3.txt",
			},
			expected: expectedPathInfo(
				"/",
				[]string{
					"home/user1/file1.txt",
					"home/user2/file2.txt",
					"var/log/file3.txt",
				},
			),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParentDirWithRelatives(tt.paths)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParentDirWithRelatives() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParentDirWithRelatives() = %+v, want %+v", got, tt.expected)
			}
		})
	}
}
func TestParentDirWithRelativesWindows(t *testing.T) {
	t.Skip("these do not work yet")
	if runtime.GOOS != "windows" {
		t.Skip("skipping test on non-Windows platforms")
	}

	tests := []struct {
		name     string
		paths    []string
		expected PathInfo
		wantErr  bool
	}{
		{
			name: "windows paths - same drive",
			paths: []string{
				"C:\\Users\\user\\docs\\file1.txt",
				"C:\\Users\\user\\docs\\subdir\\file2.txt",
			},
			expected: PathInfo{
				"C:\\Users\\user\\docs",
				[]string{"file1.txt", "subdir\\file2.txt"},
			},
			wantErr: false,
		},
		{
			name: "windows paths - different drives",
			paths: []string{
				"C:\\Users\\user\\file1.txt",
				"D:\\Documents\\file2.txt",
			},
			expected: PathInfo{
				"", // No common parent for different drives
				[]string{
					"C:\\Users\\user\\file1.txt",
					"D:\\Documents\\file2.txt",
				},
			},
			wantErr: false,
		},
		{
			name: "windows paths - UNC paths",
			paths: []string{
				"\\\\server\\share\\folder1\\file1.txt",
				"\\\\server\\share\\folder1\\folder2\\file2.txt",
			},
			expected: PathInfo{
				"\\\\server\\share\\folder1",
				[]string{"file1.txt", "folder2\\file2.txt"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParentDirWithRelatives(tt.paths)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParentDirWithRelatives() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParentDirWithRelatives() = %+v, want %+v", got, tt.expected)
			}
		})
	}
}

// TestParentDirWithRelativesWithInvalidPaths tests error handling
func TestParentDirWithRelativesWithInvalidPaths(t *testing.T) {
	t.Skip("these do not work yet")
	tests := []struct {
		name    string
		paths   []string
		wantErr bool
	}{
		{
			name: "invalid characters in path",
			paths: []string{
				string([]byte{0x00}) + "/invalid/path",
				"/valid/path",
			},
			wantErr: true,
		},
		{
			name: "malformed paths",
			paths: []string{
				"/valid/path",
				"///invalid///.././path",
			},
			wantErr: false, // Should clean the path, not error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParentDirWithRelatives(tt.paths)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParentDirWithRelatives() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
