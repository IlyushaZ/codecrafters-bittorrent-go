package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func ListPeers(ctx context.Context, addr string, info map[string]interface{}) ([]string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	be := bencoder{&buf}
	if err := be.encode(info); err != nil {
		return nil, err
	}

	h := sha1.New()
	io.Copy(h, &buf)

	q := u.Query()
	q.Add("info_hash", fmt.Sprintf("%s", h.Sum(nil)))
	q.Add("peer_id", "00112233445566778899")
	q.Add("port", "6881")
	q.Add("uploaded", "0")
	q.Add("downloaded", "0")
	q.Add("left", strconv.Itoa(info["length"].(int)))
	q.Add("compact", "1")

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bd := bdecoder{bufio.NewReader(resp.Body)}
	decoded, err := bd.decode()
	if err != nil {
		return nil, err
	}

	decodedMap, ok := decoded.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected decoded response to be map, %T given", decoded)
	}

	return parsePeerAddrs(decodedMap["peers"].(string)), nil
}

func parsePeerAddrs(src string) []string {
	var result []string

	for i := 0; i < len(src); i += 6 {
		bytesAddr := src[i : i+6]
		port := binary.BigEndian.Uint16([]byte(bytesAddr[4:6]))

		result = append(result, fmt.Sprintf("%d.%d.%d.%d:%d", bytesAddr[0], bytesAddr[1], bytesAddr[2], bytesAddr[3], port))
	}

	return result
}
