package normalizex

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func NormalizeName(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func NormalizeTitle(s string) string {
	s = strings.Join(strings.Fields(s), " ") // colapsa espacios
	lc := cases.Lower(language.Spanish)
	tc := cases.Title(language.Spanish)

	return tc.String(lc.String(s))
}

func NormalizarStreet(s string) string {
	// Trim espacios extras al inicio y final
	text := strings.TrimSpace(s)

	// Eliminar cualquier cantidad de '#'
	reHash := regexp.MustCompile(`#*`)
	text = reHash.ReplaceAllString(text, "")

	return NormalizeTitle(text)
}
