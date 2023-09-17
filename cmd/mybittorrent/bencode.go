package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"unicode"
)

var (
	ErrUnsupportedType = errors.New("unsupported bencode type")
	ErrMalformedString = errors.New("malformed string")
)

type bdecoder struct {
	*bufio.Reader
}

func (b *bdecoder) decode() (interface{}, error) {
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

type bencoder struct {
	// *bufio.Writer
	*bytes.Buffer
}

func (b *bencoder) encode(val interface{}) error {
	switch v := val.(type) {
	case string:
		b.WriteString(fmt.Sprintf("%d:%s", len(v), v))
		return nil
	case int:
		b.WriteString(fmt.Sprintf("i%de", v))
		return nil
	case []interface{}:
		b.WriteByte('l')
		for _, item := range v {
			b.encode(item)
		}
		b.WriteByte('e')

		return nil
	case map[string]interface{}:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		b.WriteByte('d')

		for _, k := range keys {
			b.encode(k)
			b.encode(v[k])
		}

		b.WriteByte('e')

		return nil
	default:
		return ErrUnsupportedType
	}
}
