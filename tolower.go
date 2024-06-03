package piscine

func ToLower(s string) string {
	tbr := []rune(s)
	for i := 0; i < len(tbr); i++ {
		if tbr[i] >= 'A' && tbr[i] <= 'Z' {
			tbr[i] += 32
		}
	}
	return string(tbr)
}
