package voice

import (
	"encoding/binary"
	"fmt"
	"golang.org/x/crypto/nacl/secretbox"
)

var silenceFrame = []byte{0xF8, 0xFF, 0xFE}

// WriteOpusFrame method takes care of all header fields, proper data encoding and writing. Will block until write is finished. Should not be used concurrently.
func (c *udpConnection) WriteOpusFrame(frame []byte) error {
	binary.BigEndian.PutUint16(c.header[2:4], c.sequence)
	c.sequence++

	c.timestamp += 960
	binary.BigEndian.PutUint32(c.header[4:8], c.timestamp)
	copy(c.nonce[:], c.header)

	writable := secretbox.Seal(c.header, frame, &c.nonce, &c.secretKey)

	select {
	case _, open := <-c.writeFrequency.C:
		if !open {
			println("not open")
			return nil
		}
	}

	if c.isClosed.Load() {
		return ConnectionClosedError
	}
	_, err := c.conn.Write(writable)
	if err != nil {
		return fmt.Errorf("error while writing packet: %w", err)
	}

	return nil
}

func (c *udpConnection) WriteSilence() {
	for i := 0; i < 5; i++ {
		_, _ = c.conn.Write(silenceFrame)
	}
}
