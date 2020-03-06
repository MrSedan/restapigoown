package apiserver

import (
	"io"
	"net/http"
	"os"
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
