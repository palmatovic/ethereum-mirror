package util

import (
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
