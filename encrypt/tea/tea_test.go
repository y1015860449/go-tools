package tea

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestEncrypt(t *testing.T) {
	type args struct {
		keyBytes   []byte
		plainBytes []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "key必须16byte", args: args{keyBytes: genRandBytes(8), plainBytes: []byte(uuid.New().String())}, wantErr: true},
		{name: "key必须16byte", args: args{keyBytes: genRandBytes(32), plainBytes: []byte(uuid.New().String())}, wantErr: true},
		{name: "解密成功", args: args{keyBytes: genRandBytes(16), plainBytes: []byte(uuid.New().String())}, wantErr: false},
		{name: "解密成功", args: args{keyBytes: genRandBytes(16), plainBytes: []byte(uuid.New().String())[:0]}, wantErr: true},
		{name: "解密成功", args: args{keyBytes: genRandBytes(16), plainBytes: []byte(uuid.New().String())[:7]}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cipherBytes, err := Encrypt(tt.args.keyBytes, tt.args.plainBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil { // 加密成功才测解密
				return
			}
			decryptBytes, err := Decrypt(tt.args.keyBytes, cipherBytes)
			if !reflect.DeepEqual(decryptBytes, tt.args.plainBytes) {
				t.Errorf("Encrypt() got = %v, want %v", decryptBytes, tt.args.plainBytes)
			}
		})
	}
}

func BenchmarkTea(b *testing.B) {
	key := GenRandTeaKey()
	var plainTxt string
	for i := 0; i < b.N; i++ {
		plainTxt = fmt.Sprintf("%s%d", plainTxt, i)
		cipherTxt, err := StrEncryptStr2Str(key, plainTxt)
		assert.Nil(b, err)
		decryptTxt, err := StrDecryptStr2Str(key, cipherTxt)
		assert.Equal(b, plainTxt, decryptTxt)
	}
}
