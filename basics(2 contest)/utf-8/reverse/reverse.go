//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	var builder strings.Builder
	reverse := make([]rune, 0, utf8.RuneCountInString(input))
	for _, r := range input {
		reverse = append([]rune{r}, reverse...)
	}
	for _, r := range reverse {
		builder.WriteRune(r)
	}
	return builder.String()
}
