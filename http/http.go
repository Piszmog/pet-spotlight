package http

import (
	"fmt"
	"net/http"
	"os"
	"pet-spotlight/io"
)

// DownloadImage downloads the file from the specified URL and saves to the provided path as the specified file
// name.
func DownloadImage(imageURL string, path string, fileName string) error {
	resp, err := http.Get(imageURL)
	if err != nil {
		return fmt.Errorf("failed to get image from %s: %w", imageURL, err)
	}
	defer io.CloseResource(resp.Body)
	imagePath := fmt.Sprintf("%s/%s", path, fileName)
	f, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", imagePath, err)
	}
	defer io.CloseResource(f)
	if err = io.CopyToFile(resp.Body, f); err != nil {
		return err
	}
	return nil
}
