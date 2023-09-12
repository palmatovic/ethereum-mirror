package util

import (
	"html"
	"strings"
)

func CleanText(text *string) {
	*text = html.UnescapeString(*text)
	*text = strings.TrimSpace(*text)
	*text = strings.TrimPrefix(*text, " ")
	*text = strings.TrimPrefix(*text, "\n")
	*text = strings.TrimPrefix(*text, "\t")
}
