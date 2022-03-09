package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const defaultTimeout = 30 * time.Second

type Client struct {
	URL    string
	Client *http.Client
}

func New(url string) *Client {
	return &Client{
		URL: url,
		Client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       defaultTimeout,
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

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	writer.Close()

	request, err := http.NewRequestWithContext(context.Background(), "POST", c.URL+"/frames", body)
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
