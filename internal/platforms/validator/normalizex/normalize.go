package normalizex

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func NormalizeTitle(s string) string {
	s = strings.Join(strings.Fields(s), " ") // colapsa espacios
	lc := cases.Lower(language.Spanish)
	tc := cases.Title(language.Spanish)

	return tc.String(lc.String(s))
}
func NormalizarStreet(input string) string {
	// Trim espacios extras al inicio y final
	text := strings.TrimSpace(input)

	// Eliminar cualquier cantidad de '#'
	reHash := regexp.MustCompile(`#*`)
	text = reHash.ReplaceAllString(text, "")

	return NormalizeTitle(text)
}
