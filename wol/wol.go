package wol

import (
	"encoding"
	"fmt"
	"net"
)

const ExpectedWrittenBytes = 102

var (
	ErrWrittenMagicPacketBytesTooLow = fmt.Errorf(
		"amount of writen bytes in the magic packet is expected to be %d",
		ExpectedWrittenBytes,
	)
)

type WOL struct {
	MagicPacket encoding.BinaryMarshaler
	Transport   net.Conn
}

func (wol *WOL) WakeUp() error {
	marshalledBytes, mpErr := wol.MagicPacket.MarshalBinary()

	if mpErr != nil {
		return mpErr
	}

	fmt.Printf("Attempting to send a magic packet to MAC %s\n", "aaa")
	writtenBytes, writeErr := wol.Transport.Write(marshalledBytes)

	if writeErr == nil && writtenBytes != ExpectedWrittenBytes {
		return fmt.Errorf(
			"wrong number of bytes sent, got %d but expected %d: %w",
			writtenBytes,
			ExpectedWrittenBytes,
			ErrWrittenMagicPacketBytesTooLow,
		)
	}

	if writeErr != nil {
		return writeErr
	}

	fmt.Println("Magic packet sent successfully")

	return wol.Transport.Close()
}

func New(magicPacket encoding.BinaryMarshaler, transp net.Conn) *WOL {
	return &WOL{
		MagicPacket: magicPacket,
		Transport:   transp,
	}
}
