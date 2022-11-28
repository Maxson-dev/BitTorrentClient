package main

import "encoding/binary"

type messageID = uint8

const (
	ChokeMessage messageID = iota
	UnchokeMessage
	InterestedMessage
	NotinterestedMessage
	HaveMessage
	BitfieldMessage
	RequestMessage
	PieceMessage
	CancelMessage
	PortMessage
)

type Message struct {
	Length  uint32
	ID      messageID
	Payload []byte
}

func (m *Message) Serialize() []byte {
	buf := make([]byte, m.Length+4)
	binary.BigEndian.PutUint32(buf[0:4], m.Length)
	offset := 4
	buf[offset] = m.ID
	offset++
	copy(buf[offset:], m.Payload)
	return buf
}

func DeserializeMessage(buf []byte) *Message {
	msg := &Message{}
	offset := 0
	msg.Length = binary.BigEndian.Uint32(buf[offset : offset+4])
	offset += 4
	msg.ID = buf[offset]
	offset++
	copy(msg.Payload, buf[offset:])
	return msg
}

func (m *Message) ParsePiece() {

}

type BitField []byte

func (b *BitField) HasPiece(index int) bool {
	if index > len(*b)*8 {
		return false
	}
	byteInd := index / 8
	num := (*b)[byteInd]
	bit := index % 8
	offset := 7 - bit
	return (num>>offset)&1 != 0
}

func (b *BitField) SetPiece(index int) {
	if index > len(*b)*8 {
		return
	}
	byteInd := index / 8
	bit := index % 8
	offset := 7 - bit
	(*b)[byteInd] |= 1 << offset
}
