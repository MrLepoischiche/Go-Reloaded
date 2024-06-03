package piscine

func Atoi(s string) int {
	if len(s) == 0 {
		return 0
	}

	sign := (s[0] == '-' || s[0] == '+')
	var negative bool
	start := 0
	if sign {
		negative = (s[0] == '-')
		start++
	}

	res := 0
	for _, char := range s[start:] {
		if char < 48 || char > 57 {
			return 0
		}
		res = res*10 + int(char-48)
	}

	if negative {
		res *= -1
	}

	return res
}
