package cronv

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

func Md5Sum(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Shorten(v string, size int, suffix string) string {
	r := []rune(v)
	if len(r) > size {
		return fmt.Sprintf("%s%s", strings.TrimSpace(string(r[:size])), suffix)
	}
	return v
}
