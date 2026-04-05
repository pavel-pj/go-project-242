package code

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Переменные для моков при тестировании больших размеров
var (
	osLstat   = os.Lstat
	osWalkDir = filepath.WalkDir // Используем WalkDir вместо Walk
)

// переменная отвечающая за форматирования при разных вызовах функции GetPathSize
var IncludePathInOutput = false

// GetPathSize возвращает общий размер всех файлов в указанном пути.
// Для директорий рекурсивно обходит все вложенные файлы и папки.
// Возвращает размер в байтах и ошибку, если путь недоступен.
func GetPathSize(path string, isRecursive, isHuman, isAll bool) (string, error) {
	size, err := getIntSize(path, isAll, isRecursive)
	if err != nil {
		return "", err
	}
	result := formatSize(size, isHuman)
	return result, nil
}

func getIntSize(path string, isAll, isRecursive bool) (uint64, error) {
	file, err := osLstat(path)
	if err != nil {
		return 0, fmt.Errorf("невозможно открыть файл : %q", path)
	}
	if !file.IsDir() {
		//nolint:gosec // dirSize всегда неотрицательный, переполнение невозможно
		return uint64(file.Size()), nil
	}

	dirSize, err := getDirSize(path, isAll, isRecursive)
	if err != nil {
		return 0, fmt.Errorf("ошибка обхода директории : %q", path)
	}
	//nolint:gosec // dirSize всегда неотрицательный, переполнение невозможно
	return uint64(dirSize), nil
}

func formatSize(size uint64, isHuman bool) string {
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
	path = strings.TrimSuffix(path, string(os.PathSeparator))
	var totalSize int64
	var walkErr error

	err := osWalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filePath == path {
			return nil
		}

		// Получаем относительный путь для проверки глубины
		relPath := strings.TrimPrefix(filePath, path)
		relPath = strings.TrimPrefix(relPath, string(os.PathSeparator))
		depth := strings.Count(relPath, string(os.PathSeparator))

		// Проверка рекурсивности
		if !isRecursive && depth > 0 {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		// Проверка на скрытые файлы
		if !d.IsDir() && (isAll || !strings.HasPrefix(d.Name(), ".")) {
			info, err := d.Info()
			if err != nil {
				// Сохраняем ошибку, но продолжаем обход
				walkErr = errors.Join(walkErr, fmt.Errorf("не удалось получить info для %s: %w", filePath, err))
				return nil
			}
			totalSize += info.Size()
		}

		return nil
	})

	if err != nil {
		return totalSize, err
	}
	if walkErr != nil {
		return totalSize, walkErr
	}
	return totalSize, nil
}
