package expense

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const serverPort = 80

func Uri(paths ...string) string {
	host := fmt.Sprintf("http://localhost:%d", serverPort)
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func Request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	token := os.Getenv("AUTH_TOKEN")
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}
