package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func gzipEncode(data string) (string, error) {
	var buffer bytes.Buffer

	writer := gzip.NewWriter(&buffer)

	_, err := writer.Write([]byte(data))

	if err != nil {
		fmt.Println("Got error at gzip encoding: ", err.Error())
	}

	writer.Close()

	return string(buffer.Bytes()), err
}
