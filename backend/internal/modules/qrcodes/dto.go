package qrcodes

import "qr-generator/backend/internal/shared"

type DesignRequest struct {
	TemplateID           *uint                       `json:"template_id"`
	ForegroundColor      string                      `json:"foreground_color"`
	BackgroundColor      string                      `json:"background_color"`
	EyeStyle             string                      `json:"eye_style"`
	DotStyle             string                      `json:"dot_style"`
	FrameStyle           string                      `json:"frame_style"`
	LogoURL              string                      `json:"logo_url"`
	Size                 int                         `json:"size"`
	ErrorCorrectionLevel shared.ErrorCorrectionLevel `json:"error_correction_level"`
}

type CreateRequest struct {
	FolderID       *uint          `json:"folder_id"`
	Title          string         `json:"title" binding:"required,max=150"`
	QRType         shared.QRType  `json:"qr_type" binding:"required"`
	Content        string         `json:"content" binding:"required"`
	IsDynamic      bool           `json:"is_dynamic"`
	DestinationURL string         `json:"destination_url"`
	Design         *DesignRequest `json:"design"`
}

type UpdateRequest struct {
	FolderID       *uint           `json:"folder_id"`
	Title          string          `json:"title" binding:"required,max=150"`
	Content        string          `json:"content"`
	DestinationURL string          `json:"destination_url"`
	Status         shared.QRStatus `json:"status"`
}
