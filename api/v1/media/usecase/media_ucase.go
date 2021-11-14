package usecase

import (
	"fmt"
	"io"
	"mime/multipart"
	"mini_pos/config"
	"os"
	"path/filepath"
	"time"

	"mini_pos/domain"
)

type MediaUsecase struct {
	cfg    config.Config
	contextTimeout time.Duration
}

// NewCustomerUsecase will create new an CustomerUsecase object representation of domain.CustomerUsecase interface
func NewCustomerUsecase(cfg config.Config, timeout time.Duration) domain.MediaUseCase {
	return &MediaUsecase{
		cfg:    cfg,
		contextTimeout: timeout,
	}
}

func (s *MediaUsecase) UploadMedia(file *multipart.FileHeader) (string, error) {
	var filename =""

	currentTime := time.Now()

	alias := currentTime.Format("20060102150405000000")+"_image"

	uploadedFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	filename = file.Filename
	if alias != "" {
		filename = fmt.Sprintf("%s%s", alias, filepath.Ext(file.Filename))
	}
	folder := "public/"
	if _, err := os.Stat(dir+"/public//"+currentTime.Format("2006-01-02")); os.IsNotExist(err) {
		err := os.Mkdir(dir+"//public//"+currentTime.Format("2006-01-02"), os.ModePerm)
		if err != nil {
			return "", err
		}
		folder = folder+"/"+currentTime.Format("2006-01-02")
	}

	fileLocation := filepath.Join(dir, folder, filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		return "", err
	}

	return "/"+currentTime.Format("2006-01-02")+"/"+filename, nil
}
