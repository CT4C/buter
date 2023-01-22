package prepare

func PreparePayloads(filenames []string) (totalPayloads int, payloadSet [][]string, err error) {
	p := make([][]string, 0)
	totalPayloads = 1

	for _, filename := range filenames {
		content, err := ReadFileByLine(filename)
		if err != nil {
			return 0, p, err
		}
		totalPayloads *= len(content)
		p = append(p, content)
	}

	return totalPayloads, p, nil
}
