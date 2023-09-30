package main

import (
	"bytes"
	"crypto/sha1"
	"io"
)

func hashInfo(info map[string]interface{}) (hash string, err error) {
	buf := bytes.Buffer{}
	be := bencoder{&buf}
	err = be.encode(info)
	if err != nil {
		return
	}

	h := sha1.New()
	io.Copy(h, &buf)

	return string(h.Sum(nil)), nil
}
