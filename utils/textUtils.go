package utils

import (
	"github.com/gosimple/slug"
)

func GenerateSlug(text1 string, text2 ...string) string {
	fullText := text1
	if len(text2) > 0 && text2[0] != "" {
		fullText += " " + text2[0]
	}
	return slug.Make(fullText)
}
