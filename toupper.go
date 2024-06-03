package piscine

func ToUpper(s string) string {
	tbr := []rune(s)
	for i := 0; i < len(tbr); i++ {
		if tbr[i] >= 'a' && tbr[i] <= 'z' {
			tbr[i] -= 32
		}
	}
	return string(tbr)
}
