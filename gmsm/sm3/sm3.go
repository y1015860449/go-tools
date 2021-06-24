package sm3

import (
	"encoding/hex"
	"github.com/tjfoc/gmsm/sm3"
)

func Sum(data []byte) string {
	rest := sm3.Sm3Sum(data)
	return hex.EncodeToString(rest)
}
