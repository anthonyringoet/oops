package files

import (
	"io/fs"
	"path/filepath"
)

type FileWalker struct {
	Root string
}

func NewFileWalker(root string) *FileWalker {
	return &FileWalker{Root: root}
}

func (f *FileWalker) WalkFiles() ([]string, error) {
	var files []string

	err := filepath.WalkDir(f.Root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
			// fmt.Printf("File processed: %s\n", path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
