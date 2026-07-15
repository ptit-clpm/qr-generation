package utils

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"
)

func GenerateShortCode() string {
	return strings.ReplaceAll(uuid.NewString()[:8], "-", "")
}

func GeneratePNG(content string, size int) ([]byte, error) {
	if size <= 0 {
		size = 512
	}
	return qrcode.Encode(content, qrcode.Medium, size)
}

func DynamicURL(appURL string, shortCode string) string {
	return fmt.Sprintf("%s/q/%s", strings.TrimRight(appURL, "/"), shortCode)
}
