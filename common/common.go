package common

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func Normalize(words []string, cutset string) []string {
	newWords := make([]string, 0, len(words))
	for _, word := range words {
		newWord := strings.ToLower(strings.Trim(word, cutset))
		if utf8.RuneCountInString(newWord) > 1 {
			newWords = append(newWords, newWord)
		}
	}
	return newWords
}

func HumanReadableSize(bytes uint64) string {
	switch {
	case bytes < 1024:
		return fmt.Sprintf("%dB", bytes)
	case bytes < 1024*1024:
		return fmt.Sprintf("%.2fK", float64(bytes)/1024)
	case bytes < 1024*1024*1024:
		return fmt.Sprintf("%.2fM", float64(bytes)/1024/1024)
	default:
		return fmt.Sprintf("%.2fG", float64(bytes)/1024/1024/1024)
	}
}
