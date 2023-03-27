package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(input []byte) string {
	h := md5.New()
	h.Write(input)
	return hex.EncodeToString(h.Sum(nil))
}
