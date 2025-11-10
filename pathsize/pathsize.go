package pathsize

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
)

// Переменные для моков при тестировании больших размеров
var (
	osLstat      = os.Lstat
	filepathWalk = filepath.Walk
)

func GetSize(path string, isHuman bool) (string, error) {

	size, err := getIntSize(path)
	if err != nil {
		return "", err
	}

	if isHuman {
		return (formatFileSize(size) + "\t" + path), nil
	}

	result := strconv.FormatUint(size, 10)
	return (result + "B\t" + path), nil

}

func getIntSize(path string) (uint64, error) {
	file, err := osLstat(path)
	if err != nil {
		return 0, fmt.Errorf("невозможно открыть файл : %q", path)
	}

	if !file.IsDir() {
		return uint64(file.Size()), nil
	}

	dirSize, err := getDirSize(path)
	if err != nil {
		return 0, fmt.Errorf("ошибка обхода директории : %q", path)
	}
	return uint64(dirSize), nil
}

func formatFileSize(size uint64) string {

	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
		PB = TB * 1024
		EB = PB * 1024
	)

	switch {
	case size > EB:
		return fmt.Sprintf("%0.1fEB", float64(size)/float64(EB))
	case size > PB:
		return fmt.Sprintf("%0.1fPB", float64(size)/float64(PB))
	case size > TB:
		return fmt.Sprintf("%0.1fTB", float64(size)/float64(TB))
	case size > GB:
		return fmt.Sprintf("%0.1fGB", float64(size)/float64(GB))
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
	err := filepathWalk(path, func(filePath string, info fs.FileInfo, err error) error {

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
