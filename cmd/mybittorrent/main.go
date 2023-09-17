package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	command := os.Args[1]

	switch command {
	case "decode":
		bencodedValue := os.Args[2]

		buf := bytes.NewBuffer([]byte(bencodedValue))
		bd := bdecoder{bufio.NewReader(buf)}
		decoded, err := bd.decode()
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))

	case "info":
		f, err := os.Open(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}

		// decode .torrent file
		bd := bdecoder{bufio.NewReader(f)}

		decoded, err := bd.decode()
		if err != nil {
			fmt.Println(err)
			return
		}

		decodedMap, ok := decoded.(map[string]interface{})
		if !ok {
			fmt.Printf("expected decoded value to be map, got %T\n", decoded)
			return
		}

		// get info section
		info := decodedMap["info"]
		infoMap, ok := info.(map[string]interface{})
		if !ok {
			fmt.Printf("expected info to be map, got %T\n", decoded)
			return
		}

		// calculate hash from info section
		buf := bytes.Buffer{}
		be := bencoder{&buf}
		err = be.encode(infoMap)
		if err != nil {
			fmt.Println(err)
			return
		}

		h := sha1.New()
		io.Copy(h, &buf)

		sum := h.Sum(nil)

		fmt.Print("Tracker URL: ", decodedMap["announce"])
		fmt.Print("Length: ", infoMap["length"])
		fmt.Printf("Info Hash: %x", sum)
		fmt.Print("Piece Length: ", infoMap["piece length"])
		fmt.Printf("Piece Hashes: %x", infoMap["pieces"])
	default:
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
