package main

import (
	"bytes"
	"crypto/sha1"
	"os"

	"github.com/jackpal/bencode-go"
)

type info struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type metaInfo struct {
	Announce string `bencode:"announce"`
	Info     info   `bencode:"info"`
}

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func parse(f string) (*metaInfo, error) {
	d, err := os.Open(f)
	defer d.Close()
	if err != nil {
		return nil, err
	}
	tor := &metaInfo{}
	err = bencode.Unmarshal(d, tor)
	if err != nil {
		return nil, err
	}
	return tor, nil
}

func getInfoHash(inf *info) ([20]byte, error) {
	buf := bytes.Buffer{}
	err := bencode.Marshal(&buf, *inf)
	if err != nil {
		return [20]byte{}, err
	}
	return sha1.Sum(buf.Bytes()), nil
}

func splitToHashes(str string) [][20]byte {
	res := make([][20]byte, len(str)/20)
	for i := 0; i < len(str); i += 20 {
		buf := bytes.Buffer{}
		var offset int
		if i+20 < len(str) {
			offset = i + 20
		} else {
			offset = len(str)
		}
		buf.WriteString(str[i:offset])
		res = append(res, sha1.Sum(buf.Bytes()))
	}
	return res
}

// NewTorrent парсит .torrent файл в структуру типа TorrentFile
func NewTorrent(f string) (*TorrentFile, error) {
	meta, err := parse(f)
	if err != nil {
		return nil, err
	}

	ins, err := getInfoHash(&meta.Info)
	if err != nil {
		return nil, err
	}

	bt := &TorrentFile{
		Announce:    meta.Announce,
		InfoHash:    ins,
		PieceHashes: splitToHashes(meta.Info.Pieces),
		PieceLength: meta.Info.PieceLength,
		Length:      meta.Info.Length,
		Name:        meta.Info.Name,
	}
	return bt, nil
}
