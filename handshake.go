package main

import (
	"bytes"
)

const (
	BitTorrentV1Pstrlen uint8  = 19
	BitTorrentV1Pstr    string = "BitTorrent protocol"
)

type Handshake struct {
	Pstrlen  uint8
	Pstr     string
	Reserved [8]byte
	InfoHash [20]byte
	PeerID   [20]byte
}

func (h *Handshake) Serialize() ([]byte, error) {
	var err error
	var buf bytes.Buffer
	buf.Grow(49 + len(h.Pstr))
	err = buf.WriteByte(h.Pstrlen)
	if err != nil {
		return nil, err
	}
	_, err = buf.WriteString(h.Pstr)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(h.Reserved[:])
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(h.InfoHash[:])
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(h.PeerID[:])
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DeserializeHandshake(buf []byte) *Handshake {
	hsh := &Handshake{}
	hsh.Pstrlen = buf[0]
	offset := 1
	pstr := make([]byte, hsh.Pstrlen)
	copy(pstr, buf[offset:hsh.Pstrlen])
	offset += int(hsh.Pstrlen)
	hsh.Pstr = b2s(pstr)
	copy(hsh.Reserved[:], buf[offset:offset+8])
	offset += 8
	copy(hsh.InfoHash[:], buf[offset:offset+20])
	offset += 20
	copy(hsh.PeerID[:], buf[offset:offset+20])
	return hsh
}
