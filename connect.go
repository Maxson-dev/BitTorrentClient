package main

import (
	"encoding/binary"
	"net"
	"strconv"
	"time"
)

type Peer struct {
	IP             net.IP
	Port           uint16
	peerChocking   bool
	peerInterested bool
}

type Connect struct {
	amChocking   bool
	amInterested bool
	Conn         *net.Conn
	*Peer
}

func NewConnect(peer *Peer, peerID [20]byte, infoHash [20]byte) (*Connect, error) {
	hsh := Handshake{
		Pstrlen:  BitTorrentV1Pstrlen,
		Pstr:     BitTorrentV1Pstr,
		Reserved: [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		InfoHash: infoHash,
		PeerID:   peerID,
	}
	buf, err := hsh.Serialize()
	if err != nil {
		return nil, err
	}
	hostPort := peer.IP.String() + ":" + strconv.Itoa(int(peer.Port))
	conn, err := net.DialTimeout("tcp", hostPort, time.Second*5)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(buf)
	if err != nil {
		return nil, err
	}
	peerConn := &Connect{
		amChocking:   true,
		amInterested: true,
		Conn:         &conn,
		Peer:         peer,
	}
	return peerConn, nil
}

func (c *Connect) SendUnchoke() error {
	msg := &Message{
		Length:  1,
		ID:      UnchokeMessage,
		Payload: []byte{},
	}
	buf := msg.Serialize()
	_, err := (*c.Conn).Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connect) SendChoke() error {
	msg := &Message{
		Length:  1,
		ID:      ChokeMessage,
		Payload: []byte{},
	}
	buf := msg.Serialize()
	_, err := (*c.Conn).Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connect) SendInterested() error {
	msg := &Message{
		Length:  1,
		ID:      InterestedMessage,
		Payload: []byte{},
	}
	buf := msg.Serialize()
	_, err := (*c.Conn).Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connect) SendNotInterested() error {
	msg := &Message{
		Length:  1,
		ID:      NotinterestedMessage,
		Payload: []byte{},
	}
	buf := msg.Serialize()
	_, err := (*c.Conn).Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connect) SendHave(pieceIndex uint32) error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, pieceIndex)
	msg := &Message{
		Length:  5,
		ID:      HaveMessage,
		Payload: buf,
	}
	data := msg.Serialize()
	_, err := (*c.Conn).Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connect) SendBitfield(bf BitField) error {
	msg := &Message{
		Length:  uint32(1 + len(bf)),
		ID:      BitfieldMessage,
		Payload: bf,
	}
	buf := msg.Serialize()
	_, err := (*c.Conn).Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connect) SendPieceRequest(pieceIndex, begin, length uint32) error {
	payload := make([]byte, 12)
	offset := 0
	binary.BigEndian.PutUint32(payload[offset:offset+4], pieceIndex)
	offset += 4
	binary.BigEndian.PutUint32(payload[offset:offset+4], begin)
	offset += 4
	binary.BigEndian.PutUint32(payload[offset:offset+4], length)
	msg := &Message{
		Length:  13,
		ID:      RequestMessage,
		Payload: payload,
	}
	buf := msg.Serialize()
	_, err := (*c.Conn).Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connect) SendPiece(pieceIndex, begin uint32, data *[]byte) error {
	payload := make([]byte, len(*data)+8)
	offset := 0
	binary.BigEndian.PutUint32(payload[offset:offset+4], pieceIndex)
	offset += 4
	binary.BigEndian.PutUint32(payload[offset:offset+4], begin)
	offset += 4
	copy(payload[offset:], *data)
	msg := &Message{
		Length:  uint32(9 + len(payload)),
		ID:      PieceMessage,
		Payload: payload,
	}
	buf := msg.Serialize()
	_, err := (*c.Conn).Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connect) SendCancel(pieceIndex, begin, length uint32) error {
	payload := make([]byte, 12)
	offset := 0
	binary.BigEndian.PutUint32(payload[offset:offset+4], pieceIndex)
	offset += 4
	binary.BigEndian.PutUint32(payload[offset:offset+4], begin)
	offset += 4
	binary.BigEndian.PutUint32(payload[offset:offset+4], length)
	msg := &Message{
		Length:  13,
		ID:      CancelMessage,
		Payload: payload,
	}
	buf := msg.Serialize()
	_, err := (*c.Conn).Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Connect) SendPort(port uint16) error {
	payload := make([]byte, 2)
	binary.BigEndian.PutUint16(payload, port)
	msg := &Message{
		Length:  3,
		ID:      PortMessage,
		Payload: payload,
	}
	buf := msg.Serialize()
	_, err := (*c.Conn).Write(buf)
	if err != nil {
		return err
	}
	return nil
}
