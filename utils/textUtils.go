package utils

import (
	"github.com/gosimple/slug"
)

func GenerateSlug(text1, text2 string) string {
	fullText := text1
	if text2 != "" {
		fullText += " " + text2
	}
	return slug.Make(fullText)
}
