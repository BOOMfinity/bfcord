package voice

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/BOOMfinity/golog"
	"go.uber.org/atomic"
	"io"
	"net"
	"time"
)

type udpConnection struct {
	conn           net.Conn
	writeFrequency *time.Ticker
	logger         golog.Logger
	isClosed       *atomic.Bool
	OwnIP          string
	header         []byte
	timestamp      uint32
	OwnPort        uint16
	sequence       uint16
	secretKey      [32]byte
	nonce          [24]byte
}

// newUDP function creates new UDP socket and performs IP discovery.
func newUDP(ctx context.Context, addr string, ssrc uint32, log golog.Logger) (*udpConnection, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	c := udpConnection{
		logger:         log.Module("UDP"),
		writeFrequency: time.NewTicker(20 * time.Millisecond),
		isClosed:       atomic.NewBool(false),
	}

	udp, err := net.Dial("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial UDP: %w", err)
	}
	c.logger.Debug().Send("UDP connected to %v with ssrc %v", addr, ssrc)
	c.conn = udp

	for i := 1; i <= 10; i++ {
		if err = c.holePunch(ssrc); err != nil {
			c.logger.Debug().Send("[%v/10] UDP holepunch failed with error: %v", i, err.Error())
			continue
		}
		c.logger.Debug().Send("[%v/10] UDP holepunch success. IP: %v:%v", i, c.OwnIP, c.OwnPort)
		break
	}
	if c.OwnIP == "" || c.OwnPort == 0 {
		return nil, UDPHolepunchFailed
	}

	c.header = make([]byte, 12)
	c.header[0], c.header[1] = 0x80, 0x78
	binary.BigEndian.PutUint32(c.header[8:12], ssrc)
	return &c, nil
}

func (c *udpConnection) PutSecretKey(key [32]byte) {
	c.secretKey = key
}

// Close method writes 5 silence frames and closes UDP socket. All calls to WriteOpusFrame after this will result in ConnectionClosedError
func (c *udpConnection) Close() {
	c.isClosed.Store(true)
	c.WriteSilence()
	_ = c.conn.Close()
}

func (c *udpConnection) holePunch(ssrc uint32) error {
	ipDiscoveryRequest := [74]byte{}
	binary.BigEndian.PutUint16(ipDiscoveryRequest[0:2], 1)
	binary.BigEndian.PutUint16(ipDiscoveryRequest[2:4], 70)
	binary.BigEndian.PutUint32(ipDiscoveryRequest[4:8], ssrc)

	_, err := c.conn.Write(ipDiscoveryRequest[:])
	if err != nil {
		return fmt.Errorf("failed to write SSRC buffer: %w", err)
	}

	var ipDiscoveryResponse = make([]byte, 74)
	rlen, err := io.ReadFull(c.conn, ipDiscoveryResponse)
	if err != nil {
		return fmt.Errorf("failed to read IP buffer: %w", err)
	}
	if rlen != 74 {
		return fmt.Errorf("IP discovery: packet length mismatch: %v received, wanted 74", rlen)
	}

	ipbody := ipDiscoveryResponse[8:72]

	nullPos := bytes.Index(ipbody, []byte{'\x00'})
	if nullPos < 0 {
		return errors.New("IP discovery response doesn't contain a null terminator")
	}

	c.OwnIP = string(ipbody[:nullPos])
	c.OwnPort = binary.LittleEndian.Uint16(ipDiscoveryResponse[72:74])

	return nil
}
