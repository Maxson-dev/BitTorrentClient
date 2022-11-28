package main

import (
	"encoding/binary"
	"fmt"
	"github.com/jackpal/bencode-go"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type trackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (t *TorrentFile) buildTrackerUrl(peerId [20]byte, port uint16) (*url.URL, error) {
	uri, err := url.Parse(t.Announce)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	v := url.Values{
		"info_hash":  []string{b2s(t.InfoHash[:])},
		"peer_id":    []string{b2s(peerId[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"compact":    []string{"1"},
		"downloaded": []string{"0"},
		"left":       []string{strconv.Itoa(t.Length)},
		"no_peer_id": []string{"1"},
	}
	uri.RawQuery = v.Encode()
	log.Println(uri.String())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return uri, nil
}

func (t *BitTorrentClient) sendTrackerRequest() (*trackerResponse, error) {
	uri, err := t.buildTrackerUrl(t.PeerID, uint16(6881))
	if err != nil {
		return nil, err
	}
	cl := http.Client{Timeout: time.Second * 15}
	resp, err := cl.Get(uri.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	tr := &trackerResponse{}
	err = bencode.Unmarshal(resp.Body, tr)
	if err != nil {
		return nil, err
	}
	return tr, nil
}

func (t *trackerResponse) extractPeers(blob string) ([]Peer, error) {
	const peerSize = 6
	buf := []byte(blob)
	numPeers := len(blob) / peerSize
	if len(blob)%peerSize != 0 {
		err := fmt.Errorf("INVALID PEERS LENGTH")
		return nil, err
	}
	peers := make([]Peer, numPeers)
	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(blob[offset : offset+4])
		fmt.Printf(peers[i].IP.String() + "\n")
		peers[i].Port = binary.BigEndian.Uint16(buf[offset+4 : offset+6])
		fmt.Printf("PORT: %d\n", peers[i].Port)
	}
	return peers, nil
}
