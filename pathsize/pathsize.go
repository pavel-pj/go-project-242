package pathsize

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func GetSize(path string) (int64, error) {

	file, err := os.Lstat(path)
	if err != nil {
		return 0, fmt.Errorf("невозможно открыть файл : %q", path)
	}

	if !file.IsDir() {
		//value := formatFileSize(file.Size())
		return file.Size(), nil
	}

	dirSize, err := getDirSize(path)
	if err != nil {
		return 0, fmt.Errorf("ошибка обхода директории : %q", path)
	}
	//dirSizeInfo := formatFileSize(dirSize)
	return dirSize, nil

}

func formatFileSize(size int64) string {

	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case size > GB:
		return fmt.Sprintf("%0.fGB", float64(size)/float64(GB))
	case size > MB:
		return fmt.Sprintf("%0.1fMB", float64(size)/float64(MB))
	case size > KB:
		return fmt.Sprintf("%0.1fKB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%dB", size)
	}
}

func getDirSize(path string) (int64, error) {

	var totalSize int64
	err := filepath.Walk(path, func(filePath string, info fs.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.IsDir() {
			totalSize += info.Size()
		}

		return nil

	})

	return totalSize, err

}
