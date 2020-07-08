package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"os"
)

func Md5SaltPass(pass string) string {
	salt := pass
	h := md5.New()
	io.WriteString(h, salt)
	io.WriteString(h, pass)
	return fmt.Sprintf("%s:%s", salt, fmt.Sprintf("%x", string(h.Sum(nil))))
}

func main() {
	var (
		password string
		help    bool
	)
	flag.StringVar(&password, "p", "", "the plain text password used to generate crypt password")
	flag.BoolVar(&help, "h", false, "help")

	flag.Usage = func() {
		fmt.Printf(`
Usage: [%s] -p password

Options:
`, os.Args[0])

		flag.PrintDefaults()
	}

	flag.Parse()

	if help || password == "" {
		flag.Usage()
		return
	}

	fmt.Println(Md5SaltPass(password))
}
