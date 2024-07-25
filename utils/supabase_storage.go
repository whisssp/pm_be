package utils

import (
	"fmt"
	storage_go "github.com/supabase-community/storage-go"
	"io"
	"mime/multipart"
	"pm/infrastructure/controllers/payload"
	"pm/infrastructure/persistences/base"
)

var supabaseStorage *storage_go.Client = nil

func InitSupabaseStorage(p *base.Persistence) {
	supabaseStorage = p.SupabaseStorage
}

func SupabaseStorageUploadFile(bucketID, fileName string, file multipart.File, contentType string) (string, error) {
	Upsert := true
	//contentType := "image/jpeg"
	_, e := uploadFileToBucket(bucketID, fileName, file, storage_go.FileOptions{
		ContentType: &contentType,
		Upsert:      &Upsert,
	})
	if e != nil {
		return "", payload.ErrDB(fmt.Errorf("error uploading file to supabase"))
	}
	urlResponse, err := SupabaseStorageGetUrl(bucketID, fileName)
	if err != nil {
		return "", payload.ErrDB(fmt.Errorf("error getting url from supabase"))
	}
	return urlResponse, nil
}

func SupabaseStorageGetUrl(bucketId, filePath string) (string, error) {
	if supabaseStorage == nil {
		return "", payload.ErrDB(fmt.Errorf("not found supabase storage client"))
	}
	url := supabaseStorage.GetPublicUrl(bucketId, filePath, storage_go.UrlOptions{
		Transform: nil,
		Download:  false,
	})
	return url.SignedURL, nil
}

func uploadFileToBucket(bucketId, filePath string, file io.Reader, fileOptions storage_go.FileOptions) (*storage_go.FileUploadResponse, error) {
	result, err := supabaseStorage.UploadFile(bucketId, filePath, file, fileOptions)
	if err != nil {
		return nil, fmt.Errorf("error uploading file to bucket - %v", err)
	}
	return &result, nil
}