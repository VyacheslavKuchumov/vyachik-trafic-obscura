package protocol

import (
	"encoding/binary"
	"fmt"
)

// Packet types
const (
	TypePing = iota
	TypePong
	TypeData
	TypeHello
)

// Packet header
type Header struct {
	Type      uint8  // Packet type
	Length    uint16 // Payload length
	Seq       uint32 // Sequence number
	Timestamp uint32 // Unix timestamp
}

// Encode header to bytes
func (h *Header) Encode() []byte {
	buf := make([]byte, 11) // 1 + 2 + 4 + 4 bytes
	buf[0] = h.Type
	binary.BigEndian.PutUint16(buf[1:3], h.Length)
	binary.BigEndian.PutUint32(buf[3:7], h.Seq)
	binary.BigEndian.PutUint32(buf[7:11], h.Timestamp)
	return buf
}

// Decode header from bytes
func (h *Header) Decode(data []byte) error {
	if len(data) < 11 {
		return fmt.Errorf("header too short")
	}
	h.Type = data[0]
	h.Length = binary.BigEndian.Uint16(data[1:3])
	h.Seq = binary.BigEndian.Uint32(data[3:7])
	h.Timestamp = binary.BigEndian.Uint32(data[7:11])
	return nil
}

// Packet represents a complete VPN packet
type Packet struct {
	Header  Header
	Payload []byte
}

// Encode entire packet
func (p *Packet) Encode() []byte {
	headerBytes := p.Header.Encode()
	return append(headerBytes, p.Payload...)
}

// Decode entire packet
func (p *Packet) Decode(data []byte) error {
	if len(data) < 11 {
		return fmt.Errorf("packet too short")
	}
	if err := p.Header.Decode(data[:11]); err != nil {
		return err
	}
	if int(p.Header.Length) != len(data)-11 {
		return fmt.Errorf("payload length mismatch")
	}
	p.Payload = data[11:]
	return nil
}
