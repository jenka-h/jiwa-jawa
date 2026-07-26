package protocol

import (
	"encoding/binary"
	"fmt"

	"jiwa-jaw/internal/error"
)

const HeaderSize = 16

type Flags uint8

const (
	FlagReliable Flags = 1 << iota
	FlagACK
	FlagEncrypted
)

func (f Flags) Has(flag Flags) bool {
	return f&flag != 0
}

type Header struct {
	Flags       Flags
	MessageType MessageType
	Reserved    uint8
	Sequence    uint32
	SessionID   uint64
}

type Packet struct {
	Header  Header
	Payload []byte
}

func NewPacket(
	sessionID uint64,
	sequence uint32,
	messageType MessageType,
	payload []byte,
	reliable bool,
) Packet {
	var flags Flags
	if reliable {
		flags |= FlagReliable
	}

	return Packet{
		Header: Header{
			Flags:       flags,
			MessageType: messageType,
			Sequence:    sequence,
			SessionID:   sessionID,
		},
		Payload: append([]byte(nil), payload...),
	}
}

func NewACK(sessionID uint64, sequence uint32) Packet {
	return Packet{
		Header: Header{
			Flags:       FlagACK,
			MessageType: MessageACK,
			Sequence:    sequence,
			SessionID:   sessionID,
		},
	}
}

func (p Packet) IsACK() bool {
	return p.Header.Flags.Has(FlagACK)
}

func (p Packet) IsReliable() bool {
	return p.Header.Flags.Has(FlagReliable)
}

// Serialization
func EncodePacket(packet Packet) ([]byte, error) {
	if !packet.Header.MessageType.IsValid() {
		return nil, fmt.Errorf(
			"%w: %s",
			err.ErrInvalidMessageType,
			packet.Header.MessageType,
		)
	}

	buf := make([]byte, HeaderSize+len(packet.Payload))

	buf[0] = byte(packet.Header.Flags)
	buf[1] = byte(packet.Header.MessageType)
	buf[2] = packet.Header.Reserved
	buf[3] = 0

	binary.BigEndian.PutUint32(buf[4:8], packet.Header.Sequence)
	binary.BigEndian.PutUint64(buf[8:16], packet.Header.SessionID)

	copy(buf[HeaderSize:], packet.Payload)

	return buf, nil
}

func DecodePacket(data []byte) (Packet, error) {
	if len(data) < HeaderSize {
		return Packet{}, fmt.Errorf(
			"%w: got %d bytes",
			err.ErrPacketTooShort,
			len(data),
		)
	}

	header := Header{
		Flags:       Flags(data[0]),
		MessageType: MessageType(data[1]),
		Reserved:    data[2],
		Sequence:    binary.BigEndian.Uint32(data[4:8]),
		SessionID:   binary.BigEndian.Uint64(data[8:16]),
	}

	if !header.MessageType.IsValid() {
		return Packet{}, fmt.Errorf(
			"%w: %s",
			err.ErrInvalidMessageType,
			header.MessageType,
		)
	}

	if header.Flags.Has(FlagACK) && header.MessageType != MessageACK {
		return Packet{}, err.ErrInvalidACK
	}

	if header.MessageType == MessageACK && !header.Flags.Has(FlagACK) {
		return Packet{}, err.ErrInvalidACK
	}

	return Packet{
		Header:  header,
		Payload: append([]byte(nil), data[HeaderSize:]...),
	}, nil
}
