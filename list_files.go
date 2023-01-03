package main

import (
	"os"
	"path/filepath"
)

func readDir(dir string) (files []string, err error) {
	err = filepath.WalkDir(dir, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}
