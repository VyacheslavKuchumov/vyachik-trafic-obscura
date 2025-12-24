package tunnel

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/VyacheslavKuchumov/vyachik-trafic-obscura/internal/protocol"
)

// ReadPacket reads a framed packet from connection
func ReadPacket(conn net.Conn) (*protocol.Packet, error) {
	// Read length prefix (2 bytes)
	lenBuf := make([]byte, 2)
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint16(lenBuf)

	// Read the actual packet
	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, err
	}

	// Decode packet
	packet := &protocol.Packet{}
	if err := packet.Decode(data); err != nil {
		return nil, err
	}

	return packet, nil
}

// WritePacket writes a framed packet to connection
func WritePacket(conn net.Conn, packet *protocol.Packet) error {
	data := packet.Encode()

	// Add length prefix
	buf := make([]byte, 2+len(data))
	binary.BigEndian.PutUint16(buf[:2], uint16(len(data)))
	copy(buf[2:], data)

	// Write to connection
	_, err := conn.Write(buf)
	return err
}

// CreateHelloPacket creates a hello packet
func CreateHelloPacket(clientID string) *protocol.Packet {
	return &protocol.Packet{
		Header: protocol.Header{
			Type:      protocol.TypeHello,
			Length:    uint16(len(clientID)),
			Seq:       0,
			Timestamp: uint32(time.Now().Unix()),
		},
		Payload: []byte(clientID),
	}
}

// CreatePingPacket creates a ping packet
func CreatePingPacket(seq uint32) *protocol.Packet {
	payload := fmt.Sprintf("PING-%d", seq)
	return &protocol.Packet{
		Header: protocol.Header{
			Type:      protocol.TypePing,
			Length:    uint16(len(payload)),
			Seq:       seq,
			Timestamp: uint32(time.Now().Unix()),
		},
		Payload: []byte(payload),
	}
}

// CreatePongPacket creates a pong packet
func CreatePongPacket(pingPacket *protocol.Packet) *protocol.Packet {
	payload := fmt.Sprintf("PONG-%d", pingPacket.Header.Seq)
	return &protocol.Packet{
		Header: protocol.Header{
			Type:      protocol.TypePong,
			Length:    uint16(len(payload)),
			Seq:       pingPacket.Header.Seq,
			Timestamp: uint32(time.Now().Unix()),
		},
		Payload: []byte(payload),
	}
}
