package files

import (
	"os"
	"path/filepath"
)

// getRootFolder returns the path to the root folder that will be used for file operations.
func getRootFolder() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// For now, for our own safety, we always return the "target" directory.
	return filepath.Join(cwd, "target"), nil
}

func GetFileList() ([]string, error) {
	root, err := getRootFolder()
	if err != nil {
		return nil, err
	}

	walker := NewFileWalker(root)
	file_list, err := walker.WalkFiles()
	if err != nil {
		return nil, err
	}

	return file_list, nil
}
