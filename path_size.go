package code

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Переменные для моков при тестировании больших размеров
var (
	osLstat      = os.Lstat
	filepathWalk = filepath.Walk
)

func GetPathSize(path string, isRecursive, isHuman, isAll bool) (string, error) {

	size, err := getIntSize(path, isAll, isRecursive)
	if err != nil {
		return "", err
	}

	result := FormatSize(size, isHuman)
	//return (result + "\t" + path), nil
	return (result), nil
}

func getIntSize(path string, isAll, isRecursive bool) (uint64, error) {

	file, err := osLstat(path)
	if err != nil {
		return 0, fmt.Errorf("невозможно открыть файл : %q", path)
	}

	if !file.IsDir() {
		return uint64(file.Size()), nil
	}

	dirSize, err := getDirSize(path, isAll, isRecursive)
	if err != nil {
		return 0, fmt.Errorf("ошибка обхода директории : %q", path)
	}
	return uint64(dirSize), nil
}

func FormatSize(size uint64, isHuman bool) string {

	if !isHuman {
		return strconv.FormatUint(size, 10) + "B"
	}

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

func getDirSize(path string, isAll, isRecursive bool) (int64, error) {

	var totalSize int64

	//Для рекурсивного обхода
	if isRecursive {
		err := filepathWalk(path, func(filePath string, info fs.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if !info.IsDir() {
				if isAll || (!isAll && !strings.HasPrefix(info.Name(), ".")) {
					totalSize += info.Size()
				}

			}

			return nil

		})

		return totalSize, err

	}

	//Для первого уровня
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	for _, entry := range entries {

		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()

		if !isAll && strings.HasPrefix(info.Name(), ".") {
			continue
		}

		if err != nil {
			return 0, err
		}
		totalSize += info.Size()

	}

	return totalSize, nil

}
