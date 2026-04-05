package code

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetSizeRealFile(t *testing.T) {

	currentDir := getTestDataPath()

	cases := []struct {
		name, path, want                      string
		isHuman, isAll, isRecursive, hasError bool
	}{
		{
			name:     "simple byte in human",
			path:     "test_B.txt",
			want:     "5B",
			isHuman:  true,
			isAll:    true,
			hasError: false,
		},
		{
			name:        "simple byte",
			path:        "test_B.txt",
			want:        "5B",
			isHuman:     false,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},
		{
			name:        "KB in human",
			path:        "test_KB.txt",
			want:        "246.7KB",
			isHuman:     true,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},
		{
			name:        "KB in bytes",
			path:        "test_KB.txt",
			want:        "252570B",
			isHuman:     false,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},

		{
			name:        "MB in human",
			path:        "file1.pdf",
			want:        "4.1MB",
			isHuman:     true,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},
		{
			name:        "MB in bytes",
			path:        "file1.pdf",
			want:        "4307732B",
			isHuman:     false,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},
		{
			name:        "MB in human",
			path:        "test_MB.pdf",
			want:        "31.9MB",
			isHuman:     true,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},
		{
			name:        "MB in bytes",
			path:        "test_MB.pdf",
			want:        "33478607B",
			isHuman:     false,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},

		{
			name:        "dir ",
			path:        "dir200",
			want:        "38038914B",
			isHuman:     false,
			isAll:       false,
			isRecursive: false,
			hasError:    false,
		},

		{
			name:        "dir -human",
			path:        "dir200",
			want:        "36.3MB",
			isHuman:     true,
			isAll:       false,
			isRecursive: false,
			hasError:    false,
		},

		{
			name:        "dir -all",
			path:        "dir200",
			want:        "71517521B",
			isHuman:     false,
			isAll:       true,
			isRecursive: false,
			hasError:    false,
		},

		{
			name:        "dir -H -all",
			path:        "dir200",
			want:        "68.2MB",
			isHuman:     true,
			isAll:       true,
			isRecursive: false,
			hasError:    false,
		},

		{
			name:        "dir -r",
			path:        "dir200",
			want:        "75825258B",
			isHuman:     false,
			isAll:       false,
			isRecursive: true,
			hasError:    false,
		},

		{
			name:        "dir -r -all",
			path:        "dir200",
			want:        "142782472B",
			isHuman:     false,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},

		{
			name:        "dir -r -H",
			path:        "dir200",
			want:        "72.3MB",
			isHuman:     true,
			isAll:       false,
			isRecursive: true,
			hasError:    false,
		},

		{
			name:        "dir -r -H -all",
			path:        "dir200",
			want:        "136.2MB",
			isHuman:     true,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},

		{
			path:     "f",
			want:     "yyyFFVDDVB",
			isHuman:  true,
			isAll:    true,
			hasError: true,
		},
	}

	for _, r := range cases {

		t.Run(r.path, func(t *testing.T) {

			path := filepath.Join(currentDir, r.path)
			got, err := GetPathSize(path, r.isRecursive, r.isHuman, r.isAll)

			if r.hasError {
				require.Error(t, err)
				require.Empty(t, got)

			} else {
				require.NoError(t, err)
				require.Equal(t, r.want, got)
			}

		})
	}

}

type mockFileInfo struct {
	name  string
	size  int64
	isDir bool
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return m.size }
func (m *mockFileInfo) Mode() fs.FileMode  { return 0644 }
func (m *mockFileInfo) ModTime() time.Time { return time.Now() }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { return nil }

// mockDirEntry реализует интерфейс fs.DirEntry для тестирования
type mockDirEntry struct {
	name  string
	size  int64
	isDir bool
}

func (m *mockDirEntry) Name() string      { return m.name }
func (m *mockDirEntry) IsDir() bool       { return m.isDir }
func (m *mockDirEntry) Type() fs.FileMode { return 0 }
func (m *mockDirEntry) Info() (fs.FileInfo, error) {
	return &mockFileInfo{
		name:  m.name,
		size:  m.size,
		isDir: m.isDir,
	}, nil
}

func TestGetSizeLargeFiles(t *testing.T) {

	// Сохраняем оригинальные функции
	originalLstat := osLstat
	originalWalkDir := osWalkDir

	// Восстанавливаем после теста
	defer func() {
		osLstat = originalLstat
		osWalkDir = originalWalkDir
	}()

	const (
		GB = 1024 * 1024 * 1024
		TB = GB * 1024
		PB = TB * 1024
		EB = PB * 1024
	)

	tests := []struct {
		name     string
		path     string
		size     int64
		isHuman  bool
		expected string
	}{

		{"GB file", "/test.gb", 3 * GB, false, "3221225472B"},
		{"GB file", "/test.gb", 3 * GB, true, "3.0GB"},

		{"TB file", "/test.tb", 2 * TB, false, "2199023255552B"},
		{"TB file", "/test.tb", 2 * TB, true, "2.0TB"},

		{"PB file", "/test.pb", 5 * PB, false, "5629499534213120B"},
		{"PB file", "/test.pb", 5 * PB, true, "5.0PB"},

		{"EB file", "/test.eb", 1 * EB, false, "1152921504606846976B"},
		{"EB file", "/test.eb", 2 * EB, true, "2.0EB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Мокаем osLstat для файлов
			osLstat = func(path string) (os.FileInfo, error) {
				return &mockFileInfo{
					name:  filepath.Base(tt.path),
					size:  tt.size,
					isDir: false,
				}, nil
			}

			// Мокаем osWalkDir для директорий (не понадобится для файлов)
			osWalkDir = func(root string, fn fs.WalkDirFunc) error {
				// Для тестов файлов эта функция не вызывается
				return nil
			}

			result, err := GetPathSize(tt.path, true, tt.isHuman, true)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected: %s, Got: %s", tt.expected, result)
			}
		})
	}
}

// TestGetSizeLargeDirectory тестирует большие директории с использованием моков
func TestGetSizeLargeDirectory(t *testing.T) {
	// Сохраняем оригинальные функции
	originalLstat := osLstat
	originalWalkDir := osWalkDir

	// Восстанавливаем после теста
	defer func() {
		osLstat = originalLstat
		osWalkDir = originalWalkDir
	}()

	const (
		GB = 1024 * 1024 * 1024
	)

	tests := []struct {
		name        string
		dirPath     string
		files       []mockDirEntry
		isRecursive bool
		isAll       bool
		isHuman     bool
		expected    string
	}{
		{
			name:    "directory with large files",
			dirPath: "/testdir",
			files: []mockDirEntry{
				{name: "file1.bin", size: 2 * GB, isDir: false},
				{name: "file2.bin", size: 3 * GB, isDir: false},
				{name: ".hidden", size: 100, isDir: false},
			},
			isRecursive: false,
			isAll:       false,
			isHuman:     false,
			expected:    "5368709120B", // 5GB в байтах
		},
		{
			name:    "directory with hidden files included",
			dirPath: "/testdir",
			files: []mockDirEntry{
				{name: "file1.bin", size: 2 * GB, isDir: false},
				{name: ".hidden", size: 100, isDir: false},
			},
			isRecursive: false,
			isAll:       true,
			isHuman:     true,
			expected:    "2.0GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Мокаем osLstat для корневой директории
			osLstat = func(path string) (os.FileInfo, error) {
				if path == tt.dirPath {
					return &mockFileInfo{
						name:  filepath.Base(tt.dirPath),
						size:  0,
						isDir: true,
					}, nil
				}
				return nil, os.ErrNotExist
			}

			// Мокаем osWalkDir для обхода директории
			osWalkDir = func(root string, fn fs.WalkDirFunc) error {
				// Сначала вызываем для корневой директории
				rootEntry := &mockDirEntry{
					name:  filepath.Base(root),
					isDir: true,
				}
				if err := fn(root, rootEntry, nil); err != nil {
					return err
				}

				// Затем для всех файлов
				for _, file := range tt.files {
					if err := fn(filepath.Join(root, file.name), &file, nil); err != nil {
						return err
					}
				}
				return nil
			}

			result, err := GetPathSize(tt.dirPath, tt.isRecursive, tt.isHuman, tt.isAll)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected: %s, Got: %s", tt.expected, result)
			}
		})
	}
}

// getTestDataPath возвращает абсолютный путь к папке testdata
func getTestDataPath() string {
	// Получаем путь к текущему файлу
	_, filename, _, _ := runtime.Caller(0)
	// Переходим в корень проекта
	projectRoot := filepath.Dir(filename)
	// Возвращаем путь к testdata
	return filepath.Join(projectRoot, "testdata")
}
