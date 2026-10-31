package lsp

import "unicode/utf16"

// byteColumn maps an LSP UTF-16 column to a byte boundary in Go source text.
func byteColumn(line string, column int) int {
	units := 0
	for offset, char := range line {
		if units >= column {
			return offset
		}
		units += utf16.RuneLen(char)
		if units > column {
			return offset
		}
	}
	return len(line)
}

func utf16Column(line string, column int) int {
	if column < 0 {
		return 0
	}
	if column > len(line) {
		column = len(line)
	}
	units := 0
	for _, char := range line[:column] {
		units += utf16.RuneLen(char)
	}
	return units
}
