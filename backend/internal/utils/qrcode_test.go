package utils

import (
	"bytes"
	"image/png"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeneratePNGWithLogoWithoutLogoKeepsQRCodeCompatible(t *testing.T) {
	pngBytes, err := GeneratePNGWithLogo("https://example.com", 512, RenderOptions{})
	require.NoError(t, err)
	image, err := png.Decode(bytes.NewReader(pngBytes))
	require.NoError(t, err)
	require.Equal(t, 512, image.Bounds().Dx())
	require.Equal(t, 512, image.Bounds().Dy())
}

func TestGeneratePNGWithLogoRejectsNonCloudinaryURL(t *testing.T) {
	_, err := GeneratePNGWithLogo("hello", 512, RenderOptions{LogoURL: "https://example.com/logo.png"})
	require.Error(t, err)
}
