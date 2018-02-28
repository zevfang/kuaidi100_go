package system

import (
	"github.com/snluu/uuid"
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
)

// 计算字符串的md5值
func Md5(source string) string {
	md5h := md5.New()
	md5h.Write([]byte(source))
	return hex.EncodeToString(md5h.Sum(nil))
}

func UUID() string {
	return uuid.Rand().Hex()
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
