package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func postFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("media", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var resultFile string = "new_" + filename
	output, err := os.Create(resultFile)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, resp.Body)
	if err != nil {
		fmt.Println("error while downloading")
		return err
	}

	return nil
}

func main() {
	target_url := "http://api.everfilter.me/filters/shinkai?nightscape=0"
	postFile(os.Args[1], target_url)
}
