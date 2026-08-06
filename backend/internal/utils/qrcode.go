package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"
	_ "golang.org/x/image/webp"
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

type RenderOptions struct {
	LogoURL    string
	Foreground string
	Background string
}

// GeneratePNGWithLogo composes a remote Cloudinary logo only after validating
// that the URL is HTTPS and belongs to Cloudinary. It never stores the image locally.
func GeneratePNGWithLogo(content string, size int, opts RenderOptions) ([]byte, error) {
	if size <= 0 {
		size = 512
	}
	if strings.TrimSpace(opts.LogoURL) == "" {
		return GeneratePNG(content, size)
	}
	parsed, err := url.Parse(opts.LogoURL)
	if err != nil || parsed.Scheme != "https" || !strings.EqualFold(parsed.Hostname(), "res.cloudinary.com") {
		return nil, fmt.Errorf("logo URL must be a Cloudinary HTTPS URL")
	}
	client := &http.Client{}
	resp, err := client.Get(parsed.String())
	if err != nil {
		return nil, fmt.Errorf("download logo: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("download logo returned status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024+1))
	if err != nil || len(data) > 5*1024*1024 {
		return nil, fmt.Errorf("logo is too large")
	}
	logo, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode logo: %w", err)
	}
	qrBytes, err := qrcode.Encode(content, qrcode.High, size)
	if err != nil {
		return nil, err
	}
	qrImage, err := png.Decode(bytes.NewReader(qrBytes))
	if err != nil {
		return nil, err
	}
	qrImage = recolor(qrImage, parseHexColor(opts.Foreground, color.RGBA{17, 24, 39, 255}), parseHexColor(opts.Background, color.RGBA{255, 255, 255, 255}))
	maxLogo := size * 18 / 100
	if maxLogo < 32 {
		maxLogo = 32
	}
	logoBounds := logo.Bounds()
	logoSize := maxLogo
	if logoBounds.Dx() > logoBounds.Dy() {
		logoSize = maxLogo * logoBounds.Dy() / logoBounds.Dx()
	}
	if logoBounds.Dy() > logoBounds.Dx() {
		logoSize = maxLogo * logoBounds.Dx() / logoBounds.Dy()
	}
	if logoSize < 1 {
		return nil, fmt.Errorf("invalid logo dimensions")
	}
	resized := image.NewRGBA(image.Rect(0, 0, logoSize, logoSize))
	draw.Draw(resized, resized.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	draw.Draw(resized, image.Rect((logoSize-logoSize*logoBounds.Dx()/max(logoBounds.Dx(), logoBounds.Dy()))/2, (logoSize-logoSize*logoBounds.Dy()/max(logoBounds.Dx(), logoBounds.Dy()))/2, (logoSize+logoSize*logoBounds.Dx()/max(logoBounds.Dx(), logoBounds.Dy()))/2, (logoSize+logoSize*logoBounds.Dy()/max(logoBounds.Dx(), logoBounds.Dy()))/2), logo, logoBounds.Min, draw.Over)
	canvas := image.NewRGBA(qrImage.Bounds())
	draw.Draw(canvas, canvas.Bounds(), qrImage, image.Point{}, draw.Src)
	x, y := (size-logoSize)/2, (size-logoSize)/2
	draw.Draw(canvas, image.Rect(x, y, x+logoSize, y+logoSize), resized, image.Point{}, draw.Over)
	var out bytes.Buffer
	if err := png.Encode(&out, canvas); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func parseHexColor(value string, fallback color.RGBA) color.RGBA {
	value = strings.TrimPrefix(strings.TrimSpace(value), "#")
	if len(value) != 6 {
		return fallback
	}
	n, err := strconv.ParseUint(value, 16, 32)
	if err != nil {
		return fallback
	}
	return color.RGBA{R: uint8(n >> 16), G: uint8(n >> 8), B: uint8(n), A: 255}
}

func recolor(src image.Image, foreground, background color.RGBA) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	for y := src.Bounds().Min.Y; y < src.Bounds().Max.Y; y++ {
		for x := src.Bounds().Min.X; x < src.Bounds().Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			if r+g+b < 3*0x8000 {
				dst.SetRGBA(x, y, foreground)
			} else {
				dst.SetRGBA(x, y, background)
			}
		}
	}
	return dst
}

func DynamicURL(appURL string, shortCode string) string {
	return fmt.Sprintf("%s/q/%s", strings.TrimRight(appURL, "/"), shortCode)
}
