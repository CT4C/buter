package prepare

import "github.com/edpryk/buter/lib/drive/files"

func PreparePayloads(filenames []string) (totalPayloads int, payloadSet [][]string, err error) {
	p := make([][]string, 0)
	totalPayloads = 1

	if len(filenames) == 0 {
		return 0, p, nil
	}

	for _, filename := range filenames {
		content, err := files.ReadByLine(filename)
		if err != nil {
			return 0, p, err
		}
		totalPayloads *= len(content)
		p = append(p, content)
	}

	return totalPayloads, p, nil
}
