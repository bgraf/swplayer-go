package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func collectFiles(searchPaths []string, fileExtensions []string) ([]string, error) {
	var files []string

	for _, path := range searchPaths {
		stat, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("collectFiles: %w", err)
		}

		if stat.Mode().IsRegular() {
			files = append(files, path)
		} else if stat.IsDir() {
			directoryFiles, err := collectFilesInDirectory(path, fileExtensions)
			if err != nil {
				return nil, fmt.Errorf("collectFiles: %w", err)
			}
			files = append(files, directoryFiles...)
		}
	}

	return files, nil
}

func collectFilesInDirectory(root string, fileExtensions []string) ([]string, error) {
	var files []string
	walkFunction := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		ext := filepath.Ext(path)
		if len(ext) == 0 {
			return nil
		}
		ext = strings.ToLower(ext[1:])

		if sliceContainsString(fileExtensions, ext) {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			files = append(files, path)
			return nil
		}
		return nil
	}
	err := filepath.Walk(root, walkFunction)
	return files, err
}

func sliceContainsString(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
