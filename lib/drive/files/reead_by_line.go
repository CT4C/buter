package files

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func ReadByLine(filename string) ([]string, error) {
	_, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filename, os.O_RDONLY, 0o400)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content := make([]string, 0)
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		content = append(content, strings.TrimSpace(line))

		if err != nil {
			return nil, err
		}
	}

	return content, nil
}
