package util

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"mime"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func FirstOrDefault[T any](slice []T, filter func(*T) bool) (element *T) {

	for i := 0; i < len(slice); i++ {
		if filter(&slice[i]) {
			return &slice[i]
		}
	}

	return nil
}

func Where[T any](slice []T, filter func(*T) bool) []*T {

	var ret []*T = make([]*T, 0)

	for i := 0; i < len(slice); i++ {
		if filter(&slice[i]) {
			ret = append(ret, &slice[i])
		}
	}

	return ret
}

func SaveBase64ToFile(base64String, filePath string) error {
	// Decode the base64 string
	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return err
	}

	// Write the decoded data to a file
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func UploadFileToMinio(ctx context.Context, filePath string) (*string, error) {
	// Load the configuration from the current directory.
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal("| cannot load config")
	}

	minioClient, err := minio.New(config.MINIO_ENDPOINT, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MINIO_ACCESS_KEY_ID, config.MINIO_SECRET_ACCESS_KEY, ""),
		Secure: config.MINIO_USE_SSL,
	})
	if err != nil {
		log.Fatalf("| Error creating MinIO client: %v", err)
	}

	// Get file name from file path
	fileName := path.Base(filePath)

	// Upload file
	info, err := minioClient.FPutObject(ctx, config.MINIO_BUCKET_NAME, fileName, filePath, minio.PutObjectOptions{ContentType: "image/jpeg"})
	if err != nil {
		log.Println("| Failed to upload to minio :", err)
		return nil, err
	}

	log.Printf("|Successfully uploaded %s of size %d to minio\n", fileName, info.Size)
	file_url := fmt.Sprintf("%s%s", config.MINIO_URL_RESULT, fileName)
	return &file_url, nil
}

func GetFileExtensionFromBase64(base64String string) (string, error) {
	// Split the base64 string to get the MIME type
	parts := strings.Split(base64String, ";")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid base64 string")
	}

	// Extract the MIME type
	mimeType := strings.TrimPrefix(parts[0], "data:")

	// Special case for image/jpeg to ensure .jpeg extension
	if mimeType == "image/jpeg" {
		return ".jpeg", nil
	}

	// Get the file extension from the MIME type
	extension, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(extension) == 0 {
		return "", fmt.Errorf("could not determine file extension")
	}

	return extension[0], nil
}

func GetBase64Data(base64String string) string {
	parts := strings.Split(base64String, ";base64,")
	return parts[1]
}

func IsValidURL(str string) bool {
	u, err := url.Parse(str)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}
