package pathsize

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSize(t *testing.T) {

	currentDir := getTestDataPath()

	cases := []struct {
		name     string
		want     int64
		hasError bool
	}{
		{
			name:     "test_B.txt",
			want:     5,
			hasError: false,
		},

		{
			name:     "test_KB.txt",
			want:     252570,
			hasError: false,
		},

		{
			name:     "file1.pdf",
			want:     4307732,
			hasError: false,
		},
		{
			name:     "test_MB.pdf",
			want:     33478607,
			hasError: false,
		},
		{
			name:     "dir200",
			want:     38038914,
			hasError: false,
		},
		{
			name:     "f",
			want:     0,
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
				require.Equal(t, r.want, got)
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
