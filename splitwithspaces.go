package piscine

func SplitWithSpaces(s string) ([]string, []rune) {
	if len(s) <= 0 {
		return []string{}, []rune{}
	}
	words := []string{}
	tmp := ""
	spaces := []rune{}

	for _, r := range s {
		if r == 0 || r == 9 || r == 10 || r == 32 {
			if tmp != "" {
				words = append(words, tmp)
				tmp = ""
				spaces = append(spaces, r)
			}
			continue
		}
		tmp += string(r)
	}
	if tmp != "" {
		words = append(words, tmp)
	}

	return words, spaces
}
