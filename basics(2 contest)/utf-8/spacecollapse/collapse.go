//go:build !solution

package spacecollapse

import "strings"

func CollapseSpaces(input string) string {
	var builder strings.Builder
	prevChar := false
	for _, r := range input {
		if r == ' ' || r == '\n' || r == '\t' || r == '\r' {
			if !prevChar {
				prevChar = true
				builder.WriteRune(' ')
			}
		} else {
			builder.WriteRune(r)
			prevChar = false
		}
	}
	return builder.String()
}
