package apiserver

import (
	"io"
	"net/http"
	"os"
	"unicode"
)

//DownloadFile downloading file from a url
func DownloadFile(fpath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func verifyPassword(s string) bool {
	number := false
	lower := false
	upper := false
	if len(s) < 8 || len(s) > 64 {
		return false
	}
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.In(c, unicode.Latin) && unicode.IsLower(c):
			lower = true
		case unicode.In(c, unicode.Latin) && unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
		default:
			return false
		}
	}
	return number && lower && upper
}

func verifyUName(s string) bool {
	if len(s) < 2 || len(s) > 15 {
		return false
	}
	if unicode.IsNumber(rune(s[0])) {
		return false
	}
	for _, c := range s {
		switch {
		case unicode.IsLetter(c) && unicode.In(c, unicode.Latin):
		case unicode.IsNumber(c):
		default:
			return false
		}
	}
	return true
}
