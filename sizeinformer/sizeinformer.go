package sizeinformer

import (
	"fmt"
	"os"
)

func GetSize(path string) int {
	file, err := os.Lstat(path)
	if err != nil {
		fmt.Printf("Невозможно открыть файл : %q", path)
	}

	return int(file.Size())

}
