package util

import (
	"math/big"
	"strings"
	"unicode"
)

func CleanText(text *string) {
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
