package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

	case first == 'd':
		result := make(map[string]interface{})
		for {
			key, err := b.decode()
			if err != nil {
				return nil, err
			}

			if key == nil {
				break
			}

			strKey, ok := key.(string)
			if !ok {
				return nil, errors.New("dictionary's key must always be string")
			}

			val, err := b.decode()
			if err != nil {
				return nil, err
			}

			if val == nil {
				break
			}

			result[strKey] = val
		}

		return result, nil

	case first == 'e':
		return nil, nil

	default:
		return nil, ErrUnsupportedType
	}
}
