package main

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func titleCase(s string) string {
	return cases.Title(language.English).String(
		strings.ReplaceAll(s, "-", " "),
	)
}
