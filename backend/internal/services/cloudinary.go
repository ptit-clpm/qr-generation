package services

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"qr-generator/backend/internal/config"

	"github.com/google/uuid"
)

var ErrCloudinaryDisabled = errors.New("Cloudinary logo upload is disabled")

type CloudinaryAsset struct {
	SecureURL string `json:"secure_url"`
	PublicID  string `json:"public_id"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Format    string `json:"format"`
	Bytes     int64  `json:"bytes"`
}

type CloudinaryService struct {
	cfg    config.Config
	client *http.Client
}

func NewCloudinaryService(cfg config.Config) *CloudinaryService {
	return &CloudinaryService{cfg: cfg, client: &http.Client{Timeout: 30 * time.Second}}
}

func (s *CloudinaryService) validate() error {
	if !s.cfg.CloudinaryEnabled {
		return ErrCloudinaryDisabled
	}
	if strings.TrimSpace(s.cfg.CloudinaryCloudName) == "" || strings.TrimSpace(s.cfg.CloudinaryAPIKey) == "" || strings.TrimSpace(s.cfg.CloudinaryAPISecret) == "" {
		return errors.New("Cloudinary is enabled but cloud name, API key, and API secret are required")
	}
	return nil
}

func (s *CloudinaryService) Upload(file io.Reader, filename string, userID, qrID uint) (CloudinaryAsset, error) {
	var asset CloudinaryAsset
	if err := s.validate(); err != nil {
		return asset, err
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	publicID := fmt.Sprintf("user_%d/qr_%d/%s", userID, qrID, uuid.NewString())
	params := map[string]string{"api_key": s.cfg.CloudinaryAPIKey, "folder": s.cfg.CloudinaryFolder, "public_id": publicID, "timestamp": timestamp}
	params["signature"] = sign(params, s.cfg.CloudinaryAPISecret)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return asset, err
	}
	if _, err = io.Copy(part, file); err != nil {
		return asset, err
	}
	for key, value := range params {
		if err = writer.WriteField(key, value); err != nil {
			return asset, err
		}
	}
	if err = writer.Close(); err != nil {
		return asset, err
	}
	endpoint := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", url.PathEscape(s.cfg.CloudinaryCloudName))
	req, err := http.NewRequest(http.MethodPost, endpoint, &body)
	if err != nil {
		return asset, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := s.client.Do(req)
	if err != nil {
		return asset, fmt.Errorf("Cloudinary upload failed: %w", err)
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return asset, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return asset, fmt.Errorf("Cloudinary upload failed with status %d", resp.StatusCode)
	}
	if err := json.Unmarshal(responseBody, &asset); err != nil {
		return asset, fmt.Errorf("invalid Cloudinary upload response: %w", err)
	}
	if !strings.HasPrefix(asset.SecureURL, "https://") || asset.PublicID == "" {
		return asset, errors.New("Cloudinary returned an invalid asset")
	}
	return asset, nil
}

func (s *CloudinaryService) Delete(publicID string) error {
	if strings.TrimSpace(publicID) == "" {
		return nil
	}
	if err := s.validate(); err != nil {
		return err
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	params := map[string]string{"api_key": s.cfg.CloudinaryAPIKey, "public_id": publicID, "timestamp": timestamp}
	params["signature"] = sign(params, s.cfg.CloudinaryAPISecret)
	form := url.Values{}
	for key, value := range params {
		form.Set(key, value)
	}
	endpoint := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/destroy", url.PathEscape(s.cfg.CloudinaryCloudName))
	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("Cloudinary delete failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Cloudinary delete failed with status %d", resp.StatusCode)
	}
	return nil
}

func sign(params map[string]string, secret string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		if key != "file" && key != "api_key" && key != "resource_type" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, key+"="+params[key])
	}
	hash := sha1.Sum([]byte(strings.Join(parts, "&") + secret))
	return hex.EncodeToString(hash[:])
}
