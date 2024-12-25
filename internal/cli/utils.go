package cli

import (
	"fmt"
	"strings"
)

// print a box around some text
func drawTextBox(txt string) string {
	// Split the text into lines
	lines := strings.Split(txt, "\n")

	// Find the width of the box (longest line + padding)
	maxLength := 0
	for _, line := range lines {
		if len(line) > maxLength {
			maxLength = len(line)
		}
	}
	width := maxLength + 4 // Add padding for the box
	// Create the top and bottom borders
	topBottomBorder := "+" + strings.Repeat("-", width-2) + "+"
	// Create the middle part with lines
	var middleLines []string
	for _, line := range lines {
		middleLines = append(middleLines, fmt.Sprintf("| %-*s |", maxLength, line))
	}
	// Join all parts together
	return fmt.Sprintf("%s\n%s\n%s", topBottomBorder, strings.Join(middleLines, "\n"), topBottomBorder)
}
