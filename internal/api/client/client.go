package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Client struct {
	URL    string
	Client *http.Client
}

func New(URL string) *Client {
	return &Client{
		URL: URL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c Client) PostFrame(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("frame", filepath.Base(file.Name()))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("POST", c.URL+"/frames", body)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

	response, err := c.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return content, nil
}
