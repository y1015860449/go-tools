package imgutils

import (
	"bytes"
	"github.com/disintegration/imaging"
	"image/gif"
	"os"
)

type ImageUtils interface {
	GetThumbData(filename string, width, height int) ([]byte, error)
	GetImageInfo(filename string) (width int, height int, length int64, isGif bool, err error)
}

type defImageUtils struct {
}

func NewImageUtils() ImageUtils {
	return &defImageUtils{}
}

// 获取缩略图
func (p *defImageUtils) GetThumbData(filename string, width, height int) ([]byte, error) {
	image, err := imaging.Open(filename)
	if err != nil {
		return nil, err
	}
	thumbImage := imaging.Resize(image, width, height, imaging.Lanczos)
	var buf bytes.Buffer
	if err = imaging.Encode(&buf, thumbImage, imaging.JPEG); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 获取图片信息
func (p *defImageUtils) GetImageInfo(filename string) (width int, height int, length int64, isGif bool, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, 0, 0, false, nil
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return 0, 0, 0, false, nil
	}
	img, err := imaging.Open(filename)
	if err != nil {
		return 0, 0, 0, false, err
	}
	_, isGifErr := gif.DecodeConfig(file)
	return img.Bounds().Max.X, img.Bounds().Max.Y, stat.Size(), isGifErr == nil, nil
}
