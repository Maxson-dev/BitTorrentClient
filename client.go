package main

import (
	"log"
)

type BitTorrentClient struct {
	PeerID [20]byte
	TorrentFile
	Peers []Peer
}

type PieceWork struct {
	PieceHash  [20]byte
	PieceIndex int
}

type PieceResult struct {
	PieceWork
	Data []byte
}

func (b *BitTorrentClient) startDownloadingFile() {

	workQueue := make(chan *PieceWork, len(b.PieceHashes))
	for i := 0; i < len(b.PieceHashes); i++ {
		task := &PieceWork{
			PieceHash:  b.PieceHashes[i],
			PieceIndex: i,
		}
		workQueue <- task
	}

	for i := 0; i < len(b.Peers); i++ {

	}
}

func (b *BitTorrentClient) spawnDownloadWorker(queue chan *PieceWork, peer *Peer) {

	conn, err := NewConnect(peer, b.PeerID, b.InfoHash)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return
	}

	for task := range queue {
		offset := 0
		buf := make([]byte, b.PieceLength)
		for {
		}
	}

}

//reader := bufio.NewReader(*conn.Conn)
//num, err := reader.Read(buf[offset:])
//offset += num
//if err != nil {
//return
//}
