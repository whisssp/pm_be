package supabase

import (
	"fmt"
	storage_go "github.com/supabase-community/storage-go"
	"pm/infrastructure/config"
)

type SupabaseStorage struct {
	StorageClient *storage_go.Client
}

func NewSupabaseStorage(appConfig *config.AppConfig) *SupabaseStorage {
	supaClient := storage_go.NewClient(appConfig.SupabaseStorageConfig.Url, appConfig.SupabaseStorageConfig.Key, nil)

	if supaClient == nil {
		fmt.Println("Failed to connect to Supabase Storage")
		return nil
	}
	fmt.Println("Connected to Supabase Storage successfully")
	return &SupabaseStorage{
		StorageClient: supaClient,
	}
}