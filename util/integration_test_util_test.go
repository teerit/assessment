//go:build unit
// +build unit

package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestUri(t *testing.T) {
	// Test URI with no path parameters
	uri := Uri()
	if uri != "http://localhost:80" {
		t.Errorf("Expected http://localhost:80, got %s", uri)
	}

	// Test URI with path parameters
	uri = Uri("expenses", "123")
	if uri != "http://localhost:80/expenses/123" {
		t.Errorf("Expected http://localhost:80/expenses/123, got %s", uri)
	}
}

func TestRequest(t *testing.T) {
	res := Request("GET", "http://example.com", nil)
	if res.Err != nil {
		t.Errorf("Expected no error, got %s", res.Err)
	}
	if res.Response.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", res.Response.StatusCode)
	}

	res = Request("GET", "invalid://url", nil)
	if res.Err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestResponse(t *testing.T) {
	res := &Response{&http.Response{StatusCode: 200}, nil}
	if res.Err != nil {
		t.Errorf("Expected no error, got %s", res.Err)
	}
	if res.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", res.StatusCode)
	}

	res = &Response{nil, fmt.Errorf("error")}
	if res.Err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestDecode(t *testing.T) {
	res := &Response{&http.Response{StatusCode: 200,
		Body: ioutil.NopCloser(strings.NewReader(`{"name":"John"}`))}, nil}
	var v struct {
		Name string `json:"name"`
	}
	err := res.Decode(&v)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	if v.Name != "John" {
		t.Errorf("Expected name John, got %s", v.Name)
	}

	res = &Response{&http.Response{StatusCode: 200,
		Body: ioutil.NopCloser(strings.NewReader(`{"name":"John`))}, nil}
	err = res.Decode(&v)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
