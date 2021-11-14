package domain

import (
	"mime/multipart"
)

// Product ...
type Media struct {
	FileName     string    `json:"filename" validate:"required"`
}


type MediaUseCase interface {
	UploadMedia(file *multipart.FileHeader) (string, error)
}
