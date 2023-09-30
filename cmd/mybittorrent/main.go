package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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

		bd := bdecoder{bufio.NewReader(f)}

		decoded, err := bd.decode()
		if err != nil {
			fmt.Println(err)
			return
		}

		decodedMap := decoded.(map[string]interface{})
		info := decodedMap["info"].(map[string]interface{})

		infoHash, err := hashInfo(info)
		if err != nil {
			fmt.Printf("can't hash info: %v", err)
			return
		}

		fmt.Print("Tracker URL: ", decodedMap["announce"])
		fmt.Print("Length: ", info["length"])
		fmt.Printf("Info Hash: %x", infoHash)
		fmt.Print("Piece Length: ", info["piece length"])
		fmt.Printf("Piece Hashes: %x", info["pieces"])

	case "peers":
		f, err := os.Open(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}

		bd := bdecoder{bufio.NewReader(f)}
		decoded, err := bd.decode()
		if err != nil {
			fmt.Println(err)
			return
		}

		decodedMap := decoded.(map[string]interface{})
		info := decodedMap["info"].(map[string]interface{})

		peers, err := ListPeers(context.Background(), decodedMap["announce"].(string), info)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Print(strings.Join(peers, ""))

	case "handshake":
		f, err := os.Open(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}

		bd := bdecoder{bufio.NewReader(f)}

		decoded, err := bd.decode()
		if err != nil {
			fmt.Println(err)
			return
		}

		decodedMap := decoded.(map[string]interface{})
		info := decodedMap["info"].(map[string]interface{})

		infoHash, err := hashInfo(info)
		if err != nil {
			fmt.Printf("can't hash info: %v", err)
			return
		}

		hs, err := Handshake(os.Args[3], "00112233445566778899", infoHash)
		if err != nil {
			fmt.Printf("can't make handshake: %v", err)
			return
		}

		fmt.Println("Peer ID:", hs.PeerID)

	default:
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
