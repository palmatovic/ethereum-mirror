package util

import (
	"strings"
)

func CleanText(text *string) {

	decodedTest := decodeText(*text)
	println(decodedTest)

	decodedTest = strings.TrimSpace(decodedTest)
	decodedTest = strings.TrimPrefix(decodedTest, " ")
	decodedTest = strings.TrimPrefix(decodedTest, "\n")
	decodedTest = strings.TrimPrefix(decodedTest, "\t")
	*text = decodedTest
}

func decodeText(text string) string {

	ascii_string = text.encode('ascii', 'ignore')
	print("Stringa in ASCII:", ascii_string.decode('ascii'))

	# Conversione da ASCII a UTF-8
	utf8_string = ascii_string.decode('ascii').encode('utf-8')
}
