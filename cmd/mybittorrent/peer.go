package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net"
)

type HandshakeResult struct {
	PeerID string
}

func Handshake(addr, peerID, infoHash string) (HandshakeResult, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return HandshakeResult{}, fmt.Errorf("can't connect to peer: %w", err)
	}
	defer conn.Close()

	var bb bytes.Buffer

	bb.WriteByte(19)
	bb.WriteString("BitTorrent protocol")
	for i := 0; i < 8; i++ {
		bb.WriteByte(0)
	}
	bb.WriteString(infoHash)
	bb.WriteString(peerID)

	if _, err := conn.Write(bb.Bytes()); err != nil {
		return HandshakeResult{}, fmt.Errorf("can't write handshake to conn: %w", err)
	}

	resp := make([]byte, bb.Len())
	if read, err := conn.Read(resp); err != nil && err != io.EOF {
		return HandshakeResult{}, fmt.Errorf("can't read response from conn: %w", err)
	} else if read != bb.Len() {
		return HandshakeResult{}, fmt.Errorf("response len %d != handshake len %d", read, bb.Len())
	}

	return HandshakeResult{PeerID: hex.EncodeToString(resp[48:])}, nil
}
