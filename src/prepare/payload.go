package prepare

func PreparePayloads(filenames []string) (variants int, payloadSet [][]string, err error) {
	p := make([][]string, 0)
	variants = 1

	for _, filename := range filenames {
		content, err := ReadFile(filename)
		if err != nil {
			return 0, p, err
		}
		variants *= len(content)
		p = append(p, content)
	}

	return variants, p, nil
}
