package validationx

import (
	"bytes"
	"fmt"
	"image"
	"net/mail"
	"strings"

	"github.com/google/uuid"
)

func IsEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidUUID(value string) bool {
	if err := uuid.Validate(value); err != nil {
		return false
	}
	return true
}

func IsPhoneNumber(value string) bool {
	return phoneRegexp.MatchString(value)
}

func IsOneTimeCode(value string) bool {
	return passwordRegexp.MatchString(value)
}

// NORMALIZE

func ClearString(value string) string {
	value = spaceRegexp.ReplaceAllString(value, " ")
	return strings.TrimSpace(value)
}

func DetectImageExtension(data []byte) (string, error) {
	nr := bytes.NewReader(data)
	_, format, err := image.Decode(nr)
	if err != nil {
		return "", fmt.Errorf("DetectImageExtension: [image decode]: [%w]", err)
	}

	switch format {
	case "jpeg":
		return "jpg", nil
	case "png":
		return "png", nil
	case "gif":
		return "gif", nil
	case "webp":
		return "webp", nil
	default:
		return "", fmt.Errorf("DetectImageExtension: [unsupported image format]: %s", format)
	}
}
