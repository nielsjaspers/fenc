package filehandling_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"os/user"

	"github.com/nielsjaspers/fenc/internal/filehandling"
)

func TestOpenFile(t *testing.T) {
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	// Create a temporary file in the home directory for testing.
	tmpFile, err := os.CreateTemp(homeDir, "openfiletest-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	relPath := strings.TrimPrefix(tmpFile.Name(), homeDir+"/")
	tests := []struct {
		name    string
		path    string
		wantErr bool
		wantPath string // for checking the opened file path
	}{
		{
			name:    "Open file with absolute path",
			path:    tmpFile.Name(),
			wantErr: false,
			wantPath: tmpFile.Name(),
		},
		{
			name:    "Open file with home tilde expansion (~/)",
			path:    "~/"+relPath,
			wantErr: false,
			wantPath: tmpFile.Name(),
		},
		{
			name:    "Open file with home only (~)",
			path:    "~",
			wantErr: false,
			wantPath: homeDir,
		},
		{
			name:    "Non-existent file returns error",
			path:    "~/this_file_does_not_exist.txt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := filehandling.OpenFile(tt.path)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("OpenFile() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("OpenFile() succeeded unexpectedly")
			}

			// Check the opened file's path (if not just ~)
			if tt.wantPath != "" {
				gotPath := ""
				if finfo, err := got.Stat(); err == nil && !finfo.IsDir() {
					// We opened a file, not a directory
					gotPathAbs, _ := filepath.Abs(got.Name())
					wantPathAbs, _ := filepath.Abs(tt.wantPath)
					gotPath = gotPathAbs
					if gotPath != wantPathAbs {
						t.Errorf("OpenFile() opened file %v, want %v", gotPath, wantPathAbs)
					}
				} else if finfo != nil && finfo.IsDir() {
					// We opened a directory (e.g. when path is ~)
					gotPath = got.Name()
					wantPathAbs, _ := filepath.Abs(tt.wantPath)
					gotPathAbs, _ := filepath.Abs(gotPath)
					if gotPathAbs != wantPathAbs {
						t.Errorf("OpenFile() opened dir %v, want %v", gotPathAbs, wantPathAbs)
					}
				}
			}
		})
	}
}
