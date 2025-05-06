package source

import (
	"context"
	"fmt"

	"io"
	"net/http"

	"github.com/IKolyas/image-previewer/internal/core/image"
)

type Storage interface {
	Get(ctx context.Context, imgData *image.ImgData) ([]byte, error)
}

func Get(ctx context.Context, imgData *image.ImgData) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", imgData.ImageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	vipsImg, err := image.NewImage(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create vips image: %w", err)
	}

	switch imgData.Action {
	case image.ImageActionFill:
		res, err := vipsImg.Fill(imgData)
		if err != nil {
			return nil, fmt.Errorf("failed to create vips image: %w", err)
		}
		return res, nil
	default:
		return nil, fmt.Errorf("action not allowed")

	}
}
