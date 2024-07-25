package application

import (
	"fmt"
	"mime/multipart"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/persistences/base"
	"pm/utils"
	"strconv"
	"time"
)

const (
	bucketID = "images"
)

type FileUsecase interface {
	UploadImage(fileName string, f multipart.File, contentType string) (string, error)
}

type fileUsecase struct {
	p *base.Persistence
}

func NewFileUsecase(p *base.Persistence) FileUsecase {
	return fileUsecase{p}
}

func (f fileUsecase) UploadImage(fileName string, file multipart.File, contentType string) (string, error) {
	instant := strconv.FormatInt(time.Now().Unix(), 10)
	result, err := utils.SupabaseStorageUploadFile(bucketID, fmt.Sprintf("%v%v", instant, fileName), file, contentType)
	if err != nil {
		fmt.Printf("Upload Failed: %v", err)
		return "", payload.ErrUploadFile(err)
	}

	return result, nil
}