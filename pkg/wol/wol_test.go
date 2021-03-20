package wol

import (
	"encoding"
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeCon struct {
	net.Conn

	write func(b []byte) (n int, err error)
	close func() error
}

func (conn *FakeCon) Write(b []byte) (n int, err error) {
	return conn.write(b)
}

func (conn *FakeCon) Close() error {
	return conn.close()
}

// interface check
var _ net.Conn = &FakeCon{}

type FakeMagicPacket struct {
	marshalBinary func() (data []byte, err error)
}

func (mp *FakeMagicPacket) MarshalBinary() (data []byte, err error) {
	return mp.marshalBinary()
}

var _ encoding.BinaryMarshaler = &FakeMagicPacket{}

func TestNewWol(t *testing.T) {
	New(
		&FakeMagicPacket{},
		&FakeCon{},
	)
}

func TestWOL_WakeUp(t *testing.T) {
	tests := []struct {
		name                     string
		magicPacketMarshalBinary func() (data []byte, err error)
		transportWrite           func(b []byte) (n int, err error)
		transporClose            func() error
		expectedResult           error
	}{
		{
			name: "HappyPath",
			magicPacketMarshalBinary: func() (data []byte, err error) {
				return []byte{}, nil
			},
			transportWrite: func(b []byte) (n int, err error) {
				return ExpectedWrittenBytes, nil
			},
			transporClose: func() error {
				return nil
			},
		},
		{
			name: "WrittenBytesErr",
			magicPacketMarshalBinary: func() (data []byte, err error) {
				return []byte{}, nil
			},
			transportWrite: func(b []byte) (n int, err error) {
				return 100, nil
			},
			transporClose: func() error {
				return nil
			},
			expectedResult: ErrWrittenMagicPacketBytesTooLow,
		},
		{
			// not setting transport overrides on purpose
			// as we should never reach that part of the code
			name: "MarshallingErr",
			magicPacketMarshalBinary: func() (data []byte, err error) {
				return []byte{}, errors.New("some error")
			},
			expectedResult: errors.New("some error"),
		},
		{
			name: "TransportWriteErr",
			magicPacketMarshalBinary: func() (data []byte, err error) {
				return []byte{}, nil
			},
			transportWrite: func(b []byte) (n int, err error) {
				return 0, errors.New("some error")
			},
			expectedResult: errors.New("some error"),
		},
		{
			name: "TransportCloseErr",
			magicPacketMarshalBinary: func() (data []byte, err error) {
				return []byte{}, nil
			},
			transportWrite: func(b []byte) (n int, err error) {
				return ExpectedWrittenBytes, nil
			},
			transporClose: func() error {
				return errors.New("some error")
			},
			expectedResult: errors.New("some error"),
		},
	}
	for _, tt := range tests {

		// this is a workaround for
		// gorountines inside loops issue
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			// this test should not use any of the unerlying implementations
			// we should only test the WOL struct
			w := New(
				&FakeMagicPacket{
					marshalBinary: tt.magicPacketMarshalBinary,
				},
				&FakeCon{
					write: tt.transportWrite,
					close: tt.transporClose,
				},
			)

			result := w.WakeUp()

			unwrapped := errors.Unwrap(result)

			if unwrapped != nil {
				assert.Equal(t, tt.expectedResult, unwrapped)
			} else {
				assert.Equal(
					t,
					tt.expectedResult,
					w.WakeUp(),
				)
			}

		})
	}
}
