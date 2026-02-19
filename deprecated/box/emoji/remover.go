package emoji

import "regexp"

// EmojiRemover Remove Emoji and make the display not shown yet.^^
func EmojiRemover(char string) string {
	emojiRegex := regexp.MustCompile(`[\x{1F600}-\x{1F64F}|[\x{1F300}-\x{1F5FF}]|[\x{1F680}-\x{1F6FF}]|[\x{1F700}-\x{1F77F}]|[\x{1F780}-\x{1F7FF}]|[\x{1F800}-\x{1F8FF}]|[\x{1F900}-\x{1F9FF}]|[\x{1FA00}-\x{1FA6F}]|[\x{1FA70}-\x{1FAFF}]|[\x{1FB00}-\x{1FBFF}]|[\x{1F170}-\x{1F251}]|[\x{1F300}-\x{1F5FF}]|[\x{1F600}-\x{1F64F}]|[\x{1FC00}-\x{1FCFF}]|[\x{1F004}-\x{1F0CF}]|[\x{1F170}-\x{1F251}]]+`)
	return emojiRegex.ReplaceAllString(char, "")
}
