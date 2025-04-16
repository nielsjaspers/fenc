package filehandling

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func OpenFile(path string) (*os.File, error) {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if path == "~" {
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}
	fmt.Printf("path: %v\n", path)

	file, err := os.Open(path)
	if err != nil {
		return &os.File{}, err
	}
	return file, nil
}
