package main

import (
	"fmt"
	"go-tools/imgutils"
)

func main() {
	testFile := "./pig.gif"
	utils := imgutils.NewImageUtils()
	width, height, length, isGif, err := utils.GetImageInfo(testFile)
	fmt.Println(width, height, length, isGif, err)
	thumbData, err := utils.GetThumbData(testFile, 10, 10)
	fmt.Println(len(thumbData), err)
}
