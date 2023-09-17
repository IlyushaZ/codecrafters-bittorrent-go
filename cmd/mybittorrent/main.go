package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
)

var (
	ErrUnsupportedType = errors.New("unsupported bencode type")
	ErrMalformedString = errors.New("malformed string")
)

type bencodeDecoder struct {
	*bufio.Reader
}

func (b *bencodeDecoder) decode() (interface{}, error) {
	first, err := b.ReadByte()
	if err != nil {
		return nil, err
	}

	switch {
	case unicode.IsDigit(rune(first)):
		lenStr, err := b.ReadString(':')
		if err != nil {
			return nil, err
		}

		lenStr = string(first) + lenStr[:len(lenStr)-1]
		// fmt.Println("len: ", lenStr)

		length, err := strconv.Atoi(lenStr)
		if err != nil {
			return nil, fmt.Errorf("can't decode length: %w", err)
		}

		s := make([]byte, length)
		read, err := b.Read(s)
		if err != nil && err != io.EOF {
			return nil, err
		}

		if read != length {
			return nil, ErrMalformedString
		}

		return string(s), nil

	case first == 'i':
		str, err := b.ReadString('e')
		if err != nil {
			return nil, err
		}

		return strconv.Atoi(str[:len(str)-1]) // exclude 'e'

	case first == 'l':
		result := []interface{}{}
		for {
			item, err := b.decode()
			if err != nil {
				return nil, err
			}

			if item == nil {
				break
			}

			result = append(result, item)
		}

		return result, nil

	case first == 'e':
		return nil, nil

	default:
		return nil, ErrUnsupportedType
	}
}

func main() {
	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]
		// fmt.Println("bencoded value: ", bencodedValue)
		buf := bytes.NewBuffer([]byte(bencodedValue))
		bd := bencodeDecoder{bufio.NewReader(buf)}

		decoded, err := bd.decode()
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
