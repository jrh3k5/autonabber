package model

import "strings"

const indentation = "  "

func buildIndentation(repetitions int) string {
	return strings.Repeat(indentation, repetitions)
}
