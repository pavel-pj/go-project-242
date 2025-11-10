package sizeinformer

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSize(t *testing.T) {

	currentDir := getTestDataPath()

	cases := []struct {
		name, want string
		hasError   bool
	}{
		{
			name:     "test_B1.txt",
			want:     fmt.Sprintf("5B\t%s/test_B.txt\n", currentDir),
			hasError: false,
		},
		{
			name:     "test_KB.txt",
			want:     fmt.Sprintf("246.7KB\t%s/test_KB.txt\n", currentDir),
			hasError: false,
		},
		{
			name:     "file1.pdf",
			want:     fmt.Sprintf("4.1MB\t%s/file1.pdf\n", currentDir),
			hasError: false,
		},
		{
			name:     "test_MB.pdf",
			want:     fmt.Sprintf("31.9MB\t%s/test_MB.pdf\n", currentDir),
			hasError: false,
		},
		{
			name:     "dir200",
			want:     fmt.Sprintf("36.3MB\t%s/dir200\n", currentDir),
			hasError: false,
		},
		{
			name:     "f",
			want:     "",
			hasError: true,
		},
	}

	for _, r := range cases {

		t.Run(r.name, func(t *testing.T) {

			path := filepath.Join(currentDir, r.name)
			got, err := GetSize(path)

			if r.hasError {
				require.Error(t, err)
				require.Empty(t, got)

			} else {
				require.NoError(t, err)
				require.Equal(t, got, r.want)
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
