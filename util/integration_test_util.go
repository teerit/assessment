package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const ServerPort = 80

func Uri(paths ...string) string {
	host := fmt.Sprintf("http://localhost:%d", ServerPort)
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func Request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	// token := os.Getenv("AUTH_TOKEN")
	token := "November 10, 2009"
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	Err error
}

func (r *Response) Decode(v interface{}) error {
	if r.Err != nil {
		return r.Err
	}

	return json.NewDecoder(r.Body).Decode(v)
}
