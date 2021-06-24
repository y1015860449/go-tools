package qrcode

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestQrCode(t *testing.T) {
	data, err := CreateQRCode("我是一个测试数据", nil)
	imageBase64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(data)
	fmt.Println(imageBase64)
	if err != nil {
		t.Fatal(err)
	}
}
