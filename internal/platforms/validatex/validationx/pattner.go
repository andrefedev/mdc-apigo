package validationx

import (
	"regexp"
)

// +1 usa
// +34 España
// +52 Mexico
// +57 Colombia
// +54 Argentina
// +56 Chile
// +51 Perú
// +58 Venezuela

const (
	spaceRegexString     = `\s+`
	alphaRegexString     = `(?i)^[A-ZAÁÉÍÑÓÚÜ]+$`
	aplhaNumRegexString  = `(?i)^[0-9A-ZÁÉÍÑÓÚÜ]+$`
	phoneRegexString     = `^3(0(0|1|2|4|5)|1\d|2[0-4]|5(0|1))\d{7}$`
	numberRegexString    = "^[0-9]+$"
	passwordRegexString  = "^[0-9]{6}"
	mapStreetRegexString = `^([0-9]+[A-Z]?)(\s*[0-9]+[A-Z]?)?$`
)

var (
	spaceRegexp     = regexp.MustCompile(spaceRegexString)
	alphaRegexp     = regexp.MustCompile(alphaRegexString)
	alphaNumRegexp  = regexp.MustCompile(aplhaNumRegexString)
	phoneRegexp     = regexp.MustCompile(phoneRegexString)
	passwordRegexp  = regexp.MustCompile(passwordRegexString)
	mapStreetRegexp = regexp.MustCompile(mapStreetRegexString)
)
