package main

import "fmt"

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
	isValid := false
	for _, e := range validEncoding {
		if e == value {
			isValid = true
			break
		}
	}
	if !isValid {
		return h
	}

	h.contentEncoding = value

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
