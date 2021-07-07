package hy_qrcode

import (
	"bytes"
	"github.com/yeqown/go-qrcode"
	"image"
)

type QROption struct {
	BgColor   string // hex background color
	FgColor   string // hex frontground color
	LogoImage image.Image
}

func CreateQRCode(text string, opts *QROption) ([]byte, error) {
	var options []qrcode.ImageOption

	if opts != nil {
		if opts.BgColor != "" {
			options = append(options, qrcode.WithBgColorRGBHex(opts.BgColor))
		}

		if opts.FgColor != "" {
			options = append(options, qrcode.WithFgColorRGBHex(opts.FgColor))
		}

		if opts.LogoImage != nil {
			options = append(options, qrcode.WithLogoImage(opts.LogoImage))
		}
	}

	qrc, err := qrcode.New(text, options...)
	if err != nil {
		return nil, err
	}
	var file bytes.Buffer
	if err := qrc.SaveTo(&file); err != nil {
		return nil, err
	}
	return file.Bytes(), nil
}
