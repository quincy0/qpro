package cryption

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encode(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
