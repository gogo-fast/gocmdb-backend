package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"
)

func Md5SaltPass(pass, salt string) string {
	if salt == "" {
		salt = RandStr(8)
	}
	h := md5.New()
	io.WriteString(h, salt)
	io.WriteString(h, pass)
	return fmt.Sprintf("%s:%s", salt, fmt.Sprintf("%x",string(h.Sum(nil))))
}

func SplitMd5SaltPass(p string) (string, string) {
	sp := strings.SplitN(p, ":", 2)
	if len(sp) < 2 {
		return "", sp[0]
	}
	return sp[0], sp[1]
}
