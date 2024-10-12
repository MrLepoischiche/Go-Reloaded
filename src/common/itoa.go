package common

func Itoa(a int) string {
	if a == 0 {
		return "0"
	}

	result := ""
	tmp := a
	if a < 0 {
		tmp = -a
	}

	for ; tmp != 0; tmp /= 10 {
		result = string(rune((tmp%10)+48)) + result
	}
	if a < 0 {
		result = "-" + result
	}

	return result
}
