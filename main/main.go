package main

import (
	"fmt"
	"os"
	"piscine-go"
)

const BIN_BASE = "01"
const DEC_BASE = "0123456789"
const HEX_BASE = "0123456789ABCDEF"

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: autocorrect [SRC] [DEST]")
		return
	}

	rfile, errOR := os.Open(os.Args[1])
	if errOR != nil {
		fmt.Fprintln(os.Stderr, "Error opening file : ", errOR.Error())
		return
	}

	bufferR := make([]byte, 4096)
	nRead, errR := rfile.Read(bufferR)
	if errR != nil && errR.Error() != "EOF" {
		fmt.Fprintln(os.Stderr, "Error reading file ", os.Args[1], " : ", errR.Error())
		rfile.Close()
		return
	}

	wfile, errOW := os.OpenFile(os.Args[2], os.O_APPEND, os.ModeAppend)
	if errOW != nil {
		fmt.Fprintln(os.Stderr, "Error opening file : ", errR.Error())
		rfile.Close()
		return
	}

	if nRead == 0 {
		fmt.Println("Nothing to read in file ", os.Args[1], ".")
	} else {
		words, spaces := piscine.SplitWithSpaces(string(bufferR))
		if spaces[len(spaces)-1] == 0 {
			spaces[len(spaces)-1] = 32
		}

		//words, spaces = autocorrect(words, spaces)
		words, spaces = format(words, spaces)
		words, spaces = autocorrect(words, spaces)

		bufferW := []byte(wordsAndSpacesJoin(words, spaces))

		nWrite, errW := wfile.Write(bufferW)
		if nWrite == 0 || errW != nil {
			fmt.Fprintln(os.Stderr, "Error writing in file ", os.Args[2], " : ", errW.Error())
			wfile.Close()
			rfile.Close()
			return
		}
	}

	errCR := rfile.Close()
	if errCR != nil {
		fmt.Fprintln(os.Stderr, "Error closing file ", os.Args[1], " : ", errCR.Error())
	}

	errCW := wfile.Close()
	if errCW != nil {
		fmt.Fprintln(os.Stderr, "Error closing file ", os.Args[2], " : ", errCW.Error())
	}
}

func format(words []string, spaces []rune) ([]string, []rune) {
	options := []string{"bin", "cap", "hex", "low", "up"}

	for wordIdx := 0; wordIdx < len(words); wordIdx++ {

		if words[wordIdx][0] == '(' && (words[wordIdx][len(words[wordIdx])-1] == ')' || words[wordIdx][len(words[wordIdx])-1] == ',') {
			fmt.Println("Formatting \"" + words[wordIdx] + "\"...")

			num := 1
			start, end := 1, 2
			for ; end < len(words[wordIdx]); end++ {
				if words[wordIdx][end] == ')' ||
					words[wordIdx][end] == '\'' ||
					((words[wordIdx][end] < 'A' || words[wordIdx][end] > 'Z') && (words[wordIdx][end] < 'a' || words[wordIdx][end] > 'z')) {
					break
				}
			}
			end--

			if indexOf(options, piscine.ToLower(words[wordIdx][start:end+1])) == -1 {
				continue
			}

			tagEndIdx := end + 1
			isNextWordTag := false

			if words[wordIdx][tagEndIdx] == ',' {
				fmt.Println("Number found...")
				for ; tagEndIdx < len(words[wordIdx]) && words[wordIdx][tagEndIdx] != ')'; tagEndIdx++ {
				}

				num = piscine.Atoi(words[wordIdx][end+2 : tagEndIdx])
				if num == 0 {
					tagEndIdx = 0
					for ; tagEndIdx < len(words[wordIdx+1]) && words[wordIdx+1][tagEndIdx] != ')'; tagEndIdx++ {
					}
					num = piscine.Atoi(words[wordIdx+1][:tagEndIdx])

					if num == 0 {
						words[wordIdx+1] = words[wordIdx+1][tagEndIdx+1:]
						if words[wordIdx+1] == "" {
							words = removeAt(words, wordIdx+1)
						}
						spaces = removeAt(spaces, wordIdx)
						words = removeAt(words, wordIdx)
						spaces = removeAt(spaces, wordIdx-1)

						continue
					} else {
						isNextWordTag = true
					}
				}
				fmt.Println("It's", num, ".")
			}

			for wordProcessIdx := 0; wordProcessIdx < num && wordIdx-wordProcessIdx > 0; wordProcessIdx++ {
				switch piscine.ToLower(words[wordIdx][start : end+1]) {
				case "bin":
					num, err := piscine.AtoiBase(words[wordIdx-wordProcessIdx-1], BIN_BASE)
					if err == nil {
						words[wordIdx-wordProcessIdx-1] = piscine.Itoa(num)
					}

				case "cap":
					fmt.Println("Capitalizing", words[wordIdx-wordProcessIdx-1], "...")
					words[wordIdx-wordProcessIdx-1] = piscine.Capitalize(words[wordIdx-wordProcessIdx-1])

				case "hex":
					num, err := piscine.AtoiBase(piscine.ToUpper(words[wordIdx-wordProcessIdx-1]), HEX_BASE)
					if err == nil {
						words[wordIdx-wordProcessIdx-1] = piscine.Itoa(num)
					}

				case "low":
					fmt.Println("Lowering", words[wordIdx-wordProcessIdx-1], "...")
					words[wordIdx-wordProcessIdx-1] = piscine.ToLower(words[wordIdx-wordProcessIdx-1])
					continue

				case "up":
					fmt.Println("Uppering", words[wordIdx-wordProcessIdx-1], "...")
					words[wordIdx-wordProcessIdx-1] = piscine.ToUpper(words[wordIdx-wordProcessIdx-1])
				}
			}

			if isNextWordTag {
				words[wordIdx+1] = words[wordIdx+1][tagEndIdx+1:]
				if words[wordIdx+1] == "" {
					words = removeAt(words, wordIdx+1)
					spaces = removeAt(spaces, wordIdx)
				}
				words = removeAt(words, wordIdx)
				spaces = removeAt(spaces, wordIdx-1)
			} else {
				words[wordIdx] = words[wordIdx][:start-1] + words[wordIdx][tagEndIdx+1:]
				if words[wordIdx] == "" {
					words = removeAt(words, wordIdx)
					spaces = removeAt(spaces, wordIdx-1)
				}
			}

			wordIdx--
		}
	}

	return words, spaces
}

func autocorrect(words []string, spaces []rune) ([]string, []rune) {
	anChars := []rune{'A', 'E', 'H', 'I', 'O', 'U', 'a', 'e', 'h', 'i', 'o', 'u'}
	punctuation := []rune{'.', '!', '?', ',', ':', ';'}

	isInsideSQuotes, isInsideDQuotes := false, false

	// First, autocorrecting...
	for wordIdx := 0; wordIdx < len(words); wordIdx++ { // For each word...

		fmt.Println("Autocorrecting", words[wordIdx], "...")

		if len(words[wordIdx]) == 1 || piscine.ToLower(words[wordIdx]) == "an" {
			switch words[wordIdx] { // One switch for all trivial cases
			case "a", "A":
				for charIdx := 0; charIdx < len(words[wordIdx+1]); charIdx++ { // "a" or "A" detected. Is it appropriate according to next word?
					char := words[wordIdx+1][charIdx]
					if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') { // Is first *letter* of next word a vowel OR an "h"?
						if indexOf(anChars, rune(char)) != -1 {
							words[wordIdx] += "n"
						} else {
							break
						}
					}
				}

			case "an", "An", "aN", "AN":
				for charIdx := 0; charIdx < len(words[wordIdx+1]); charIdx++ { // "an" or "An" or "aN" or "AN" detected. Is it appropriate according to next word?
					char := words[wordIdx+1][charIdx]
					if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') { // Is first *letter* of next word NOT a vowel AND NOT an "h"?
						if indexOf(anChars, rune(char)) == -1 {
							words[wordIdx] = words[wordIdx][:1]
						} else {
							break
						}
					}
				}

			case "(", "[", "{":
				words[wordIdx+1] = words[wordIdx] + words[wordIdx+1] // Current word duplicated at start of next word
				words = removeAt(words, wordIdx)                     // Current word is deleted
				spaces = removeAt(spaces, wordIdx)                   // Space AFTER current word is removed.
				wordIdx--

			case ")", "]", "}":
				words[wordIdx-1] = words[wordIdx-1] + words[wordIdx] // Current word duplicated at end of previous word
				words = removeAt(words, wordIdx)                     // Current word is deleted
				spaces = removeAt(spaces, wordIdx-1)                 // Space BEFORE current word is removed.
				wordIdx--

			case "'", "\"":
				switch words[wordIdx] {
				case "\"":
					if isInsideDQuotes {
						words[wordIdx-1] = words[wordIdx-1] + words[wordIdx] // Current word duplicated at end of previous word
						words = removeAt(words, wordIdx)                     // Current word is deleted
						spaces = removeAt(spaces, wordIdx-1)                 // Space BEFORE current word is removed.
					} else {
						words[wordIdx+1] = words[wordIdx] + words[wordIdx+1] // Current word duplicated at start of next word
						words = removeAt(words, wordIdx)                     // Current word is deleted
						spaces = removeAt(spaces, wordIdx)                   // Space AFTER current word is removed.
					}
					isInsideDQuotes = !isInsideDQuotes // Toggle boolean value, as state has changed
				case "'":
					if isInsideSQuotes {
						words[wordIdx-1] = words[wordIdx-1] + words[wordIdx] // Current word duplicated at end of previous word
						words = removeAt(words, wordIdx)                     // Current word is deleted
						spaces = removeAt(spaces, wordIdx-1)                 // Space BEFORE current word is removed.
					} else {
						words[wordIdx+1] = words[wordIdx] + words[wordIdx+1] // Current word duplicated at start of next word
						words = removeAt(words, wordIdx)                     // Current word is deleted
						spaces = removeAt(spaces, wordIdx)                   // Space AFTER current word is removed.
					}
					isInsideSQuotes = !isInsideSQuotes // Toggle boolean value, as state has changed
				}
				wordIdx--

			case ".", "!", "?", ",", ":", ";":
				if wordIdx > 0 {
					switch words[wordIdx] {
					case ".":
						prdCnt := 0
						for ; prdCnt < 3 && words[wordIdx-1][len(words[wordIdx-1])-(prdCnt+1)] == '.'; prdCnt++ { // How many periods at end of previous word?
						}
						if prdCnt < 3 {
							words[wordIdx-1] = words[wordIdx-1] + words[wordIdx] // Current word duplicated at end of previous word
							words = removeAt(words, wordIdx)                     // Current word is deleted
							spaces = removeAt(spaces, wordIdx-1)                 // Space BEFORE current word is removed.
						}
						wordIdx--
					case "?", "!":
						if indexOf(punctuation, rune((words[wordIdx-1])[len(words[wordIdx-1])-1])) <= 2 { // Any punctuation allowing groups at end of previous word?
							words[wordIdx-1] = words[wordIdx-1] + words[wordIdx] // Current word duplicated at end of previous word
							words = removeAt(words, wordIdx)                     // Current word is deleted
							spaces = removeAt(spaces, wordIdx-1)                 // Space BEFORE current word is removed.
						}
					case ",", ":", ";":
						if indexOf(punctuation, rune((words[wordIdx-1])[len(words[wordIdx-1])-1])) != 1 { // NO punctuation at end of previous word?
							words[wordIdx-1] = words[wordIdx-1] + words[wordIdx] // Current word duplicated at end of previous word
							words = removeAt(words, wordIdx)                     // Current word is deleted
							spaces = removeAt(spaces, wordIdx-1)                 // Space BEFORE current word is removed.
						}
					}
				}
			}
			continue
		}

		// Search for modifications inside word
		for charIdx := 0; charIdx < len(words[wordIdx]); charIdx++ { // For each character of current word...

			switch words[wordIdx][charIdx] {
			case '(', '[', '{':
				if charIdx > 0 {
					if charIdx == len(words[wordIdx])-1 {
						words[wordIdx+1] = words[wordIdx][charIdx:] + words[wordIdx+1]
						words[wordIdx] = words[wordIdx][:charIdx]
						continue
					} else {
						words = insertAt(words, words[wordIdx][:charIdx], wordIdx+1)
						spaces = insertAt(spaces, ' ', wordIdx+1)
						words[wordIdx] = words[wordIdx][charIdx:]
					}
				}
			case ')', ']', '}', '.', '!', '?', ',', ':', ';':
				if charIdx < len(words[wordIdx])-1 {
					if charIdx == 0 {
						switch words[wordIdx][charIdx] {
						case '.':
							prdCnt := 0
							for ; prdCnt < 3 && words[wordIdx-1][len(words[wordIdx-1])-(prdCnt+1)] == '.'; prdCnt++ { // How many periods at end of previous word?
							}
							if prdCnt < 3 {
								words[wordIdx-1] = words[wordIdx-1] + words[wordIdx][:charIdx+1] // Current character duplicated at end of previous word
								words[wordIdx] = words[wordIdx][charIdx+1:]                      // Current character is deleted
							}

						case '!', '?':
							switch (words[wordIdx-1])[len(words[wordIdx-1])-1] {
							case '.', '!', '?':
								words[wordIdx-1] = words[wordIdx-1] + words[wordIdx][:charIdx+1] // Current character duplicated at end of previous word
								words[wordIdx] = words[wordIdx][charIdx+1:]                      // Current character is deleted
							default:
								words = insertAt(words, words[wordIdx][:charIdx+1], wordIdx)
								spaces = insertAt(spaces, ' ', wordIdx)
								words[wordIdx] = words[wordIdx][charIdx+1:]
							}

						default:
							switch (words[wordIdx-1])[len(words[wordIdx-1])-1] {
							case '.', '!', '?', ',', ':', ';':
								words = insertAt(words, words[wordIdx][:charIdx+1], wordIdx)
								spaces = insertAt(spaces, ' ', wordIdx)
								words[wordIdx] = words[wordIdx][charIdx+1:]
							default:
								words[wordIdx-1] = words[wordIdx-1] + words[wordIdx][:charIdx]
								words[wordIdx] = words[wordIdx][charIdx+1:]
							}
						}
						charIdx--
						continue
					} else {
						switch words[wordIdx][charIdx] {
						case ')', ']', '}', '.', '!', '?':
							endOfWord := charIdx + 1
							for ; endOfWord < len(words[wordIdx]); endOfWord++ {
								if words[wordIdx][endOfWord] != ')' && words[wordIdx][endOfWord] != ']' && words[wordIdx][endOfWord] != '}' && words[wordIdx][endOfWord] != '\'' && words[wordIdx][endOfWord] != '"' {
									break
								}
							}
							if endOfWord < len(words[wordIdx])-1 {
								words = insertAt(words, words[wordIdx][endOfWord:], wordIdx+1)
								spaces = insertAt(spaces, ' ', wordIdx+1)
								words[wordIdx] = words[wordIdx][:endOfWord]
							}
						default:
							words = insertAt(words, words[wordIdx][charIdx+1:], wordIdx+1)
							spaces = insertAt(spaces, ' ', wordIdx+1)
							words[wordIdx] = words[wordIdx][:charIdx+1]
						}
						break
					}
				}
			case '"':
				if isInsideDQuotes && charIdx < len(words[wordIdx])-1 {
					if charIdx == 0 {
						words[wordIdx-1] = words[wordIdx-1] + words[wordIdx][:charIdx]
						words[wordIdx] = words[wordIdx][charIdx+1:]
						charIdx--
					} else {
						words = insertAt(words, words[wordIdx][:charIdx+1], wordIdx)
						spaces = insertAt(spaces, ' ', wordIdx)
						words[wordIdx] = words[wordIdx][charIdx+1:]
					}
				}
				if !isInsideDQuotes && charIdx > 0 {
					if charIdx == len(words[wordIdx])-1 {
						words[wordIdx+1] = words[wordIdx][charIdx:] + words[wordIdx+1]
						words[wordIdx] = words[wordIdx][:charIdx]
					} else {
						words = insertAt(words, words[wordIdx][:charIdx], wordIdx+1)
						spaces = insertAt(spaces, ' ', wordIdx+1)
						words[wordIdx] = words[wordIdx][charIdx:]
					}
				}
				isInsideDQuotes = !isInsideDQuotes

			case '\'':
				if charIdx != len(words[wordIdx])-2 || (words[wordIdx][len(words[wordIdx])-1] != 'M' && words[wordIdx][len(words[wordIdx])-1] != 'S' && words[wordIdx][len(words[wordIdx])-1] != 'T' && words[wordIdx][len(words[wordIdx])-1] != 'm' && words[wordIdx][len(words[wordIdx])-1] != 's' && words[wordIdx][len(words[wordIdx])-1] != 't') {
					if isInsideSQuotes && charIdx < len(words[wordIdx])-1 {
						if charIdx == 0 {
							words[wordIdx-1] = words[wordIdx-1] + words[wordIdx][:charIdx]
							words[wordIdx] = words[wordIdx][charIdx+1:]
							charIdx--
						} else {
							words = insertAt(words, words[wordIdx][:charIdx+1], wordIdx)
							spaces = insertAt(spaces, ' ', wordIdx)
							words[wordIdx] = words[wordIdx][charIdx+1:]
						}
					}
					if !isInsideSQuotes && charIdx > 0 {
						if charIdx == len(words[wordIdx])-1 {
							words[wordIdx+1] = words[wordIdx][charIdx:] + words[wordIdx+1]
							words[wordIdx] = words[wordIdx][:charIdx]
						} else {
							words = insertAt(words, words[wordIdx][:charIdx], wordIdx+1)
							spaces = insertAt(spaces, ' ', wordIdx+1)
							words[wordIdx] = words[wordIdx][charIdx:]
						}
					}
					isInsideSQuotes = !isInsideSQuotes
				}
			}
		}
	}

	if len(spaces) > len(words)-1 {
		for i := len(spaces) - 1; i > 0 && len(spaces) > len(words)-1; i-- {
			if spaces[i] == '\n' {
				continue
			}
			spaces = removeAt(spaces, i)
		}
	}

	return words, spaces
}

func indexOf[T comparable](arr []T, match T) int {
	if len(arr) <= 0 {
		return -1
	}

	for i, elem := range arr {
		if elem == match {
			return i
		}
	}

	return -1
}

func removeAt[T any](arr []T, idx int) []T {
	result := make([]T, 0)

	for i, elem := range arr {
		if i == idx {
			continue
		}
		result = append(result, elem)
	}

	return result
}

func insertAt[T any](arr []T, elem T, idx int) []T {
	result := make([]T, len(arr)+1)

	for i, where := 0, 0; where < len(result); i, where = i+1, where+1 {
		if i == idx {
			result[where] = elem
			where++
		}
		result[where] = arr[i]
	}

	return result
}

func wordsAndSpacesJoin(words []string, spaces []rune) string {
	if len(words) == 0 {
		return ""
	}
	if len(spaces) == 0 {
		return basicJoin(words)
	}

	result := ""

	for idxW, idxS := 0, 0; idxW < len(words) || idxS < len(spaces); idxW, idxS = idxW+1, idxS+1 {
		if idxW < len(words) {
			result += words[idxW]
		}
		if idxS < len(spaces) {
			result += string(spaces[idxS])
		}
	}

	return result
}

func basicJoin[T string | rune](arr []T) string {
	if len(arr) == 0 {
		return ""
	}

	result := ""

	for _, elem := range arr {
		result += string(elem)
	}

	return result
}
