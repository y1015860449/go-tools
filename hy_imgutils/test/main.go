package main

import (
	"fmt"
	"github.com/y1015860449/go-tools/hy_imgutils"
)

func main() {
	testFile := "./pig.gif"
	utils := hy_imgutils.NewImageUtils()
	width, height, length, isGif, err := utils.GetImageInfo(testFile)
	fmt.Println(width, height, length, isGif, err)
	thumbData, err := utils.GetThumbData(testFile, 10, 10)
	fmt.Println(len(thumbData), err)
}
