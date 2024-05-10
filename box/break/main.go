package Stringbreaker

import (
	"strings"
	"unicode/utf8"
)

var UserFloatNum float64

// BreakWords Break words into pieces and make adjust here.
func BreakWords(breakWords string, LimitedLength float64) string {
	var setBreaker bool
	var charCount = 0.0
	var truncated string
	for _, runeValue := range breakWords {
		charWidth := utf8.RuneLen(runeValue)
		if charWidth != 3 {
			UserFloatNum = 1.5
		} else {
			UserFloatNum = float64(charWidth)
		}
		if charCount+UserFloatNum > LimitedLength {
			setBreaker = true
			break
		}
		truncated += string(runeValue)
		charCount += UserFloatNum
	}
	if setBreaker {
		return truncated + "..."
	} else {
		return truncated
	}
}

// SplitCommandTo Split Command and Adjust To.
func SplitCommandTo(raw string, setCommandStopper int) (splitCommandLen int, splitInfo []string) {
	rawSplit := strings.SplitN(raw, " ", setCommandStopper)
	return len(rawSplit), rawSplit
}
