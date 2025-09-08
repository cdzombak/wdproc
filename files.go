package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type FileInfo struct {
	Path        string
	Name        string
	NameNoExt   string
	Content     string
	CreatedTime time.Time
	Age         time.Duration
}

func ScanDirectory(dirPath string, minAge time.Duration, now time.Time) ([]FileInfo, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var files []FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		filePath := filepath.Join(dirPath, name)
		fileInfo, err := processFile(filePath, name, minAge, now)
		if err != nil {
			return nil, fmt.Errorf("failed to process file %s: %w", filePath, err)
		}

		if fileInfo != nil {
			files = append(files, *fileInfo)
		}
	}

	return files, nil
}

func processFile(filePath, name string, minAge time.Duration, now time.Time) (*FileInfo, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	if !stat.Mode().IsRegular() {
		return nil, nil
	}

	createdTime, err := getFileCreationTime(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get creation time: %w", err)
	}

	age := now.Sub(createdTime)
	if age < minAge {
		return nil, nil
	}

	content, err := readFileContent(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	nameNoExt := strings.TrimSuffix(name, filepath.Ext(name))

	return &FileInfo{
		Path:        filePath,
		Name:        name,
		NameNoExt:   nameNoExt,
		Content:     content,
		CreatedTime: createdTime,
		Age:         age,
	}, nil
}

func readFileContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func getFileCreationTime(filePath string) (time.Time, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}

	sys := stat.Sys().(*syscall.Stat_t)
	
	createdTime := time.Unix(sys.Birthtimespec.Sec, sys.Birthtimespec.Nsec)
	return createdTime, nil
}