package docs

type PayloadFiles []string

func (ps *PayloadFiles) Set(value string) error {
	*ps = append(*ps, value)
	return nil
}

func (ps *PayloadFiles) String() string {
	s := ""
	for _, v := range *ps {
		s += v
	}
	return s
}
