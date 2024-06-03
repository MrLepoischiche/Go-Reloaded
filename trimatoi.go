package piscine

func TrimAtoi(s string) int {
	res := 0
	negative := false
	for _, r := range s {
		if r == '-' && res == 0 {
			negative = true
		}
		if r >= '0' && r <= '9' {
			res = res*10 + int(r-48)
		}
	}
	if negative {
		res = -res
	}
	return res
}
