package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// MD5 获取文件的md5, 转换为大写
func MD5(file []byte) string {
	h := md5.New()
	h.Write(file)
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
