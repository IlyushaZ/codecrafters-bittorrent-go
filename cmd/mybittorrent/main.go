package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	command := os.Args[1]

	switch command {
	case "decode":
		bencodedValue := os.Args[2]
		buf := bytes.NewBuffer([]byte(bencodedValue))
		bd := bencodeDecoder{bufio.NewReader(buf)}

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

		bd := bencodeDecoder{bufio.NewReader(f)}

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

		info := decodedMap["info"]
		infoMap, ok := info.(map[string]interface{})
		if !ok {
			fmt.Printf("expected info to be map, got %T\n", decoded)
			return
		}

		fmt.Print("Tracker URL: ", decodedMap["announce"])
		fmt.Print("Length: ", infoMap["length"])
	default:
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
