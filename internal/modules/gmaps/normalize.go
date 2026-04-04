package gmaps

import (
	"strings"

	"apigo/internal/platforms/validatex/normalizex"
)

func isValidLatLng(lat, lng float64) bool {
	if lat < -90 || lat > 90 {
		return false
	}
	if lng < -180 || lng > 180 {
		return false
	}
	return !(lat == 0 && lng == 0)
}

func normalizeStreet(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	return normalizex.NormalizarStreet(value)
}

func normalizeTitleOrKeep(value string) string {
	value = normalizeQuery(value)
	if value == "" {
		return ""
	}

	return normalizex.NormalizeTitle(value)
}

func isApproximateLocationType(value string) bool {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "APPROXIMATE":
		return true
	case "GEOMETRIC_CENTER":
		return true
	default:
		return false
	}
}

func normalizeQuery(input string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(input)), " ")
}
