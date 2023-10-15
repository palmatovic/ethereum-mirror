package wallet_transaction

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
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

	weirdSubscriptFormat, err := convertFloatAndLetter(weirdSubscriptFormat)
	if err != nil {
		return 0, err
	}
	expanded := subscripts.ReplaceAllStringFunc(weirdSubscriptFormat, func(s string) string {
		//fmt.Println("Replacing:", s)
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

	returnFloat, err := strconv.ParseFloat(expanded, 64)
	if err != nil {
		return 0, err
	}
	return returnFloat, nil
}

func convertFloatAndLetter(input string) (string, error) {
	lastChar := input[len(input)-1:]
	switch lastChar {
	case "B":
		deRuner(&input)
		stringWithoutLastCharacter := input[:len(input)-1]
		floatValue, err := strconv.ParseFloat(stringWithoutLastCharacter, 64)
		if err != nil {
			logrus.WithField("input", input).WithError(err).Errorf("cannot parse value: %s", stringWithoutLastCharacter)
		}
		return fmt.Sprintf("%2.f", floatValue*math.Pow(10, 9)), nil
	case "Q":
		deRuner(&input)
		stringWithoutLastCharacter := input[:len(input)-1]
		floatValue, err := strconv.ParseFloat(stringWithoutLastCharacter, 64)
		if err != nil {
			logrus.WithField("input", input).WithError(err).Errorf("cannot parse value: %s", stringWithoutLastCharacter)
		}
		return fmt.Sprintf("%2.f", floatValue*math.Pow(10, 15)), nil
	case "M":
		deRuner(&input)
		stringWithoutLastCharacter := input[:len(input)-1]
		floatValue, err := strconv.ParseFloat(stringWithoutLastCharacter, 64)
		if err != nil {
			logrus.WithField("input", input).WithError(err).Errorf("cannot parse value: %s", stringWithoutLastCharacter)
		}
		return fmt.Sprintf("%2.f", floatValue*math.Pow(10, 6)), nil
	default:
		return input, nil
	}
}

func deRuner(text *string) {
	newText := ""
	for _, r := range *text {
		parsed := runeToAscii(r)
		newText = newText + parsed
	}
	*text = newText
}

func runeToAscii(r rune) string {
	if r < 128 {
		return string(r)
	} else {
		asciiString := "\\u" + strconv.FormatInt(int64(r), 16)
		char, _ := strconv.Unquote(asciiString)
		return char
	}
}
