package utils

import (
	"strings"
)

/*
TrimString trims a string to a certain width and adds an ellipsis if the string is longer than the width.
The ellipsis is placed in the middle of the string, with the prefixChars and suffixChars characters before and after the ellipsis.
*/
func TrimString(s string, width int, prefixChars int, suffixChars int, ellipsis string) string {
	if len(s) > width {
		ellipsisLength := len(ellipsis)
		keepLength := width - ellipsisLength
		if keepLength <= 0 {
			return ellipsis
		}

		startLength := prefixChars
		if startLength > keepLength {
			startLength = keepLength / 2
		}
		endLength := keepLength - startLength

		return s[:startLength] + ellipsis + s[len(s)-endLength:]
	}
	return s
}

func CenterString(s string, width int) string {
	if len(s) >= width {
		return s
	}
	spaces := width - len(s)
	leftPadding := spaces / 2
	rightPadding := spaces - leftPadding
	return strings.Repeat(" ", leftPadding) + s + strings.Repeat(" ", rightPadding)
}
