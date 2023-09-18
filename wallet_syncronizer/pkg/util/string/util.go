package string

import (
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func CleanText(text *string) {
	//*text = strings.Map(func(r rune) rune {
	//	if unicode.IsPrint(r) {
	//		return r
	//	}
	//	return -1
	//}, *text)
	//*text = strings.ToValidUTF8(*text, "")
	*text = strings.TrimSpace(*text)
	*text = strings.TrimPrefix(*text, " ")
	*text = strings.TrimPrefix(*text, "\n")
	*text = strings.TrimPrefix(*text, "\t")
}
func CleanTextWithRemoveUnicodeSpaces(text *string) {
	*text = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, *text)
	*text = strings.ToValidUTF8(*text, "")
	*text = strings.TrimSpace(*text)
	*text = strings.TrimPrefix(*text, " ")
	*text = strings.TrimPrefix(*text, "\n")
	*text = strings.TrimPrefix(*text, "\t")
}

func EmptyTokenBalance(tokenBalance string) bool {
	return tokenBalance == "0x0000000000000000000000000000000000000000000000000000000000000000"
}

func CalculateAmount(tokenHexAmount string, decimals int) float64 {
	intValue, _ := new(big.Int).SetString(tokenHexAmount, 0)
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	floatValue := new(big.Float).Quo(new(big.Float).SetInt(intValue), new(big.Float).SetInt(scale))

	s, _ := floatValue.Float64()
	return s
}

const (
	subscript0Rune = rune(0x2080)
	subscript9Rune = rune(0x2089)
)

var subscripts = regexp.MustCompile(`\d\p{Zs}*[₀-₉]+\p{Zs}*`)

func ParseScript(weirdSubscriptFormat string) (float64, error) {
	if isValidFloatAndLetter(weirdSubscriptFormat) {
		return convertFloatAndLetter(weirdSubscriptFormat)
	}
	expanded := subscripts.ReplaceAllStringFunc(weirdSubscriptFormat, func(s string) string {
		fmt.Println("Replacing:", s)
		toRepeat := s[0:1]
		repeatCount := 0
		for _, rune := range s[1:] {
			if rune >= subscript0Rune && rune <= subscript9Rune {
				repeatCount *= 10
				repeatCount += int(rune - subscript0Rune)
			}
		}
		return strings.Repeat(toRepeat, repeatCount)
	})
	fmt.Println("converted to:", expanded)
	return strconv.ParseFloat(expanded, 64)
}

func isValidFloatAndLetter(input string) bool {
	// Regular expression to match a float and a letter
	re := regexp.MustCompile(`(\d+\.\d*[a-zA-Z]|[a-zA-Z]\d+\.\d*|\d*\.\d+[a-zA-Z]|[a-zA-Z]\d*\.\d*)`)
	return re.MatchString(input)
}

func convertFloatAndLetter(input string) (float64, error) {
	if isValidFloatAndLetter(input) {
		// Find the letter in the input
		var letter rune
		for _, char := range input {
			if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' {
				letter = char
				break
			}
		}

		// Extract the numeric part of the input
		numericPart := strings.ReplaceAll(input, string(letter), "")
		numericValue, err := strconv.ParseFloat(numericPart, 64)
		if err != nil {
			return 0, err
		}

		// Determine the multiplier based on the letter
		multiplier := 1.0
		switch letter {
		case 'B', 'b':
			multiplier = 1e9
		case 'Q', 'q':
			multiplier = 1e15
			// Add more cases for other letters as needed
		}

		result := numericValue * multiplier
		return result, nil
	}

	return 0, fmt.Errorf("invalid input %s", input)
}
