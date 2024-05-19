package main

import (
	"fmt"
	"strings"
)

const CLRF = "\r\n"

type Headers struct {
	contentEncoding string
	contentType     string
	contentLength   int
}

func NewHeaders() *Headers {
	return &Headers{
		contentLength: -1,
	}
}

func (h *Headers) setContentType(value string) *Headers {
	h.contentType = value
	return h
}

func (h *Headers) setContentLength(value int) *Headers {
	h.contentLength = value
	return h
}

var validEncoding = [...]string{
	"gzip",
}

func (h *Headers) setEncoding(value string) *Headers {
	values := strings.Split(value, ",")

	found := ""
	for _, value := range values {
		for _, enc := range validEncoding {
			if enc == strings.TrimSpace(value) {
				found = value
				break
			}
		}
	}
	if found == "" {
		return h
	}

	h.contentEncoding = found

	return h
}

func (h *Headers) String() string {
	headers := ""

	if h.contentEncoding != "" {
		headers +=
			fmt.Sprintf("Content-Encoding: %s", h.contentEncoding) + CLRF
	}

	if h.contentType != "" {
		headers +=
			fmt.Sprintf("Content-Type: %s", h.contentType) + CLRF
	}

	if h.contentLength > -1 {
		headers +=
			fmt.Sprintf("Content-Length: %d", h.contentLength) + CLRF
	}

	return headers
}
