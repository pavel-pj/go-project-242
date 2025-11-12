package pathsize

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
			want:     "5B\t" + filepath.Join(currentDir, "test_B.txt"),
			isHuman:  true,
			isAll:    true,
			hasError: false,
		},
		{
			name:        "simple byte",
			path:        "test_B.txt",
			want:        "5B\t" + filepath.Join(currentDir, "test_B.txt"),
			isHuman:     false,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},
		{
			name:        "KB in hyman",
			path:        "test_KB.txt",
			want:        "246.7KB\t" + filepath.Join(currentDir, "test_KB.txt"),
			isHuman:     true,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},
		{
			name:        "KB in bytes",
			path:        "test_KB.txt",
			want:        "252570B\t" + filepath.Join(currentDir, "test_KB.txt"),
			isHuman:     false,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},

		{
			name:        "MB in human",
			path:        "file1.pdf",
			want:        "4.1MB\t" + filepath.Join(currentDir, "file1.pdf"),
			isHuman:     true,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},
		{
			name:        "MB in bytes",
			path:        "file1.pdf",
			want:        "4307732B\t" + filepath.Join(currentDir, "file1.pdf"),
			isHuman:     false,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},
		{
			name:        "MB in human",
			path:        "test_MB.pdf",
			want:        "31.9MB\t" + filepath.Join(currentDir, "test_MB.pdf"),
			isHuman:     true,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},
		{
			name:        "MB in bytes",
			path:        "test_MB.pdf",
			want:        "33478607B\t" + filepath.Join(currentDir, "test_MB.pdf"),
			isHuman:     false,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},

		{
			name:        "dir ",
			path:        "dir200",
			want:        "38038914B\t" + filepath.Join(currentDir, "dir200"),
			isHuman:     false,
			isAll:       false,
			isRecursive: false,
			hasError:    false,
		},

		{
			name:        "dir -human",
			path:        "dir200",
			want:        "36.3MB\t" + filepath.Join(currentDir, "dir200"),
			isHuman:     true,
			isAll:       false,
			isRecursive: false,
			hasError:    false,
		},

		{
			name:        "dir  -all",
			path:        "dir200",
			want:        "71517521B\t" + filepath.Join(currentDir, "dir200"),
			isHuman:     false,
			isAll:       true,
			isRecursive: false,
			hasError:    false,
		},

		{
			name:        "dir -H -all",
			path:        "dir200",
			want:        "68.2MB\t" + filepath.Join(currentDir, "dir200"),
			isHuman:     true,
			isAll:       true,
			isRecursive: false,
			hasError:    false,
		},

		{
			name:        "dir -r",
			path:        "dir200",
			want:        "75825258B\t" + filepath.Join(currentDir, "dir200"),
			isHuman:     false,
			isAll:       false,
			isRecursive: true,
			hasError:    false,
		},

		{
			name:        "dir -r -all",
			path:        "dir200",
			want:        "142782472B\t" + filepath.Join(currentDir, "dir200"),
			isHuman:     false,
			isAll:       true,
			isRecursive: true,
			hasError:    false,
		},

		{
			name:        "dir -r -H",
			path:        "dir200",
			want:        "72.3MB\t" + filepath.Join(currentDir, "dir200"),
			isHuman:     true,
			isAll:       false,
			isRecursive: true,
			hasError:    false,
		},

		{
			name:        "dir -r -H -all",
			path:        "dir200",
			want:        "136.2MB\t" + filepath.Join(currentDir, "dir200"),
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
			got, err := GetSize(path, r.isHuman, r.isAll, r.isRecursive)

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

func TestGetSizeLargeFiles(t *testing.T) {

	// Сохраняем оригинальные функции
	originalLstat := osLstat
	originalWalk := filepathWalk

	// Восстанавливаем после теста
	defer func() {
		osLstat = originalLstat
		filepathWalk = originalWalk
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

		{"GB file", "/test.gb", 3 * GB, false, "3221225472B\t/test.gb"},
		{"GB file", "/test.gb", 3 * GB, true, "3.0GB\t/test.gb"},

		{"TB file", "/test.tb", 2 * TB, false, "2199023255552B\t/test.tb"},
		{"TB file", "/test.tb", 2 * TB, true, "2.0TB\t/test.tb"},

		{"PB file", "/test.pb", 5 * PB, false, "5629499534213120B\t/test.pb"},
		{"PB file", "/test.pb", 5 * PB, true, "5.0PB\t/test.pb"},

		{"EB file", "/test.eb", 1 * EB, false, "1152921504606846976B\t/test.eb"},
		{"EB file", "/test.eb", 2 * EB, true, "2.0EB\t/test.eb"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			osLstat = func(path string) (os.FileInfo, error) {
				return &mockFileInfo{
					name:  filepath.Base(tt.path),
					size:  tt.size,
					isDir: false,
				}, nil
			}

			result, err := GetSize(tt.path, tt.isHuman, true, true)
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

	// Получаем путь к текущему файлу (sizeinformer_test.go)
	_, filename, _, _ := runtime.Caller(0)

	// Переходим в корень проекта (на два уровня выше от sizeinformer_test.go)
	projectRoot := filepath.Dir(filepath.Dir(filename))

	// Возвращаем путь к testdata
	return filepath.Join(projectRoot, "testdata")
}
