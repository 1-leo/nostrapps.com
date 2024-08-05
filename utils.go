package main

import (
	"encoding/json"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func titleCase(s string) string {
	return cases.Title(language.English).String(
		strings.ReplaceAll(s, "-", " "),
	)
}

func jsEncode(value any) string {
	j, _ := json.Marshal(value)
	return string(j)
}
