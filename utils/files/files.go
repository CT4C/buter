package files

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// func ReadFiles(files []string) ([][]string, error) {
// 	payloads := make([][]string, len(files))

// 	for _, filename := range files {
// 		content, err := ReadFile(filename)
// 		if err != nil {
// 			return nil, err
// 		}
// 		payloads = append(payloads, content)
// 	}

// 	return payloads, nil
// }

func ReadFileByLine(filename string) ([]string, error) {

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
		content = append(content, strings.TrimSpace(line))

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}
	}

	return content, nil
}
