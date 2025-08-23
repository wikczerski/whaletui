package shared

// TruncName safely truncates a string to a maximum length for display purposes.
// If the string is longer than maxLen, it returns the first maxLen characters.
// If the string is shorter or equal to maxLen, it returns the original string.
func TruncName(name string, maxLen int) string {
	if len(name) > maxLen {
		return name[:maxLen]
	}
	return name
}
