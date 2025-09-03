package utils

import (
	"strings"
)

// TruncateText truncates text to the specified length and adds ellipsis if needed
func TruncateText(text string, maxLength int) string {
	if maxLength <= 0 {
		return text
	}

	if len(text) <= maxLength {
		return text
	}

	// Reserve space for ellipsis
	if maxLength <= 3 {
		return text[:maxLength]
	}

	return text[:maxLength-3] + "..."
}

// TruncateTextWithEllipsis truncates text to the specified length with custom ellipsis
func TruncateTextWithEllipsis(text string, maxLength int, ellipsis string) string {
	if maxLength <= 0 {
		return text
	}

	if len(text) <= maxLength {
		return text
	}

	ellipsisLen := len(ellipsis)
	if maxLength <= ellipsisLen {
		return text[:maxLength]
	}

	return text[:maxLength-ellipsisLen] + ellipsis
}

// TruncateTextSmart truncates text but tries to break at word boundaries when possible
func TruncateTextSmart(text string, maxLength int) string {
	if maxLength <= 0 {
		return text
	}

	if len(text) <= maxLength {
		return text
	}

	// Reserve space for ellipsis
	if maxLength <= 3 {
		return text[:maxLength]
	}

	truncateAt := maxLength - 3

	// Try to find a word boundary (space) near the truncation point
	if truncateAt > 0 {
		// Look for the last space within the last 10 characters before truncation
		searchStart := truncateAt - 10
		if searchStart < 0 {
			searchStart = 0
		}

		lastSpace := strings.LastIndex(text[searchStart:truncateAt], " ")
		if lastSpace != -1 {
			// Found a space, truncate there
			truncateAt = searchStart + lastSpace
		}
	}

	return text[:truncateAt] + "..."
}
