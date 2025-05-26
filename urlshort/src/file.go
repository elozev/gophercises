package urlshort

import (
	"fmt"
	"io"
	"os"
)

func ParseFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)

	defer file.Close()

	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", filename)
	}

	return io.ReadAll(file)
}
