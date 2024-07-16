//go:build !solution

package speller

import (
	"strings"
)

var unitNames = []string{
	"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
}

var teenNames = []string{
	"ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen",
}

var tensNames = []string{
	"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety",
}

var hundredsNames = "hundred"

var powersOfThousand = []string{
	"", "thousand", "million", "billion",
}

func spellThreeDigits(n int) string {
	if n == 0 {
		return ""
	}

	hundreds := n / 100
	remainder := n % 100
	result := ""

	if hundreds > 0 {
		result += unitNames[hundreds] + " " + hundredsNames
	}

	if remainder > 0 {
		if result != "" {
			result += " "
		}

		if remainder < 10 {
			result += unitNames[remainder]
		} else if remainder < 20 {
			result += teenNames[remainder-10]
		} else {
			tens := remainder / 10
			units := remainder % 10
			result += tensNames[tens]
			if units > 0 {
				result += "-" + unitNames[units]
			}
		}
	}

	return result
}

func Spell(n int64) string {
	if n == 0 {
		return "zero"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	blocks := []int{}
	for n > 0 {
		blocks = append(blocks, int(n%1000))
		n /= 1000
	}

	parts := []string{}
	for i := 0; i < len(blocks); i++ {
		if blocks[i] > 0 {
			spelling := spellThreeDigits(blocks[i])
			if powersOfThousand[i] != "" {
				spelling += " " + powersOfThousand[i]
			}
			parts = append([]string{spelling}, parts...)
		}
	}

	result := strings.Join(parts, " ")

	if negative {
		result = "minus " + result
	}

	return result
}
