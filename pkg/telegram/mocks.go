package telegram

import (
	"io"
	"net/http"
	"strings"
)

type BodyMock struct {
	Reader io.Reader
}

func (bm BodyMock) Read(p []byte) (n int, err error) {
	n, err = bm.Reader.Read(p)
	return
}

func (bm BodyMock) Close() error {
	return nil
}

type httpClientMock struct {
	Body string
}

func (client httpClientMock) Do(httpRequest *http.Request) (*http.Response, error) {
	httpResponse := http.Response{}
	httpResponse.Request = httpRequest
	httpResponse.Body = BodyMock{strings.NewReader(client.Body)}
	return &httpResponse, nil
}
