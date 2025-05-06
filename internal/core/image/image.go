package image

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/davidbyttow/govips/v2/vips"
)

type ImageAction string

const (
	ImageActionFill ImageAction = "fill"
)

type ImageInterface interface {
	Resize(width, height int) error
	Thumbnail(width, height int) error
	Export(format vips.ImageType) ([]byte, error)
	Fill()
	Convert()
}

type ImgData struct {
	ImageURL string
	Width    int
	Height   int
	Format   vips.ImageType
	Action   ImageAction
}

func (img *ImgData) String() string {

	v := reflect.Indirect(reflect.ValueOf(img))

	var ss []string
	for i := range v.NumField() {
		ss = append(ss, fmt.Sprintf("%v", v.Field(i).Interface()))
	}

	return strings.Join(ss, "|")
}

type Image struct {
	VipsImg *vips.ImageRef
	ImageInterface
}

func NewImage(imgData []byte) (*Image, error) {
	img, err := vips.NewImageFromBuffer(imgData)
	if err != nil {
		return nil, fmt.Errorf("failed to load image: %w", err)
	}

	return &Image{VipsImg: img}, nil
}

func (i *Image) Fill(imgData *ImgData) ([]byte, error) {
	err := i.resize(imgData.Width, imgData.Height)
	if err != nil {
		return nil, fmt.Errorf("failed to resize image: %w", err)
	}

	err = i.thumbnail(imgData.Width, imgData.Height)
	if err != nil {
		return nil, fmt.Errorf("failed to thumbnail image: %w", err)
	}

	result, err := i.export()
	if err != nil {
		return nil, fmt.Errorf("failed to export image: %w", err)
	}

	return result, nil
}

func (i *Image) resize(width, height int) error {
	scale := calculateScale(i.VipsImg, width, height)
	err := i.VipsImg.Resize(scale, vips.KernelLanczos3)
	if err != nil {
		return fmt.Errorf("failed to resize image: %w", err)
	}
	return nil
}

func (i *Image) thumbnail(width, height int) error {
	if width > 0 && height > 0 {
		err := i.VipsImg.Thumbnail(width, height, vips.InterestingCentre)
		if err != nil {
			return fmt.Errorf("failed to thumbnail image: %w", err)
		}
	}
	return nil
}

func (i *Image) export() ([]byte, error) {

	format := i.VipsImg.Metadata().Format

	switch format {
	case vips.ImageTypeJPEG:
		params := vips.NewJpegExportParams()
		params.Quality = 85
		params.OptimizeCoding = true
		imageBytes, _, err := i.VipsImg.ExportJpeg(params)
		if err != nil {
			return nil, fmt.Errorf("failed to export image: %w", err)
		}
		return imageBytes, nil

	case vips.ImageTypePNG:
		params := vips.NewPngExportParams()
		params.Compression = 6
		params.Interlace = false
		imageBytes, _, err := i.VipsImg.ExportPng(params)
		if err != nil {
			return nil, fmt.Errorf("failed to export image: %w", err)
		}
		return imageBytes, nil

	case vips.ImageTypeWEBP:
		params := vips.NewWebpExportParams()
		params.Quality = 80
		params.Lossless = false
		params.ReductionEffort = 4
		imageBytes, _, err := i.VipsImg.ExportWebp(params)
		if err != nil {
			return nil, fmt.Errorf("failed to export image: %w", err)
		}
		return imageBytes, nil

	default:
		return nil, fmt.Errorf("failed to export image. Unknown format")
	}
}

// вычисляет коэффициент масштабирования с сохранением пропорций
func calculateScale(img *vips.ImageRef, width, height int) float64 {
	if width == 0 && height == 0 {
		return 1.0
	}

	imgWidth := float64(img.Width())
	imgHeight := float64(img.Height())

	if width > 0 && height > 0 {
		// Если заданы оба размера, выбираем минимальный масштаб
		scaleW := float64(width) / imgWidth
		scaleH := float64(height) / imgHeight
		return min(scaleW, scaleH)
	} else if width > 0 {
		return float64(width) / imgWidth
	} else {
		return float64(height) / imgHeight
	}
}
