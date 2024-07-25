package utils

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"strings"
	"unicode"
)

func removeAccents(input string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, input)
	return result
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

func SanitizeFileName(fileName string) string {
	fileName = removeAccents(fileName)
	fileName = strings.ReplaceAll(fileName, " ", "_")
	fileName = strings.ReplaceAll(fileName, "/", "_")
	fileName = strings.ReplaceAll(fileName, "\\", "_")
	return fileName
}