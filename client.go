// Package modbusclient implements Modbus TCP client.
package modbusclient

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

// Errors:
var (
	ErrNotConnected  = errors.New("not connected")
	ErrShortResponse = errors.New("short response")
	ErrModbusError   = errors.New("modbus error")
	ErrShortPayload  = errors.New("short payload")
)

// Modbus Parameters:
const (
	readFuncCode     byte = 0x03
	readResHeaderLen      = 9
	errCodeIndex          = 7
	excCodeIndex          = 8
)

// Client defines the behaviors of a Modbus TCP Client.
type Client interface {
	// Connect uses net.DialTimeout to establish an underlying TCP connection with the Modbus TCP server.
	Connect() error

	// SetDeadline sets the underlying TCP connection's deadline. Returns a modbusclient.ErrNotconnected if the client is not connected.
	SetDeadline(t time.Time) error

	// Read reads data from the Holding Registers of a Modbus TCP device and writes it to the provided payload. Returns the read-byte count and a modbusclient.ErrNotconnected if the client is not connected to the server.
	Read(p []byte, unitID byte, addr uint16, count uint16) (n int, err error)

	// ReadErr parses and returns the Modbus read error of the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.
	ReadErr(p []byte) error

	// Uint16 parses and returns an uint16 value from the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.
	Uint16(p []byte, offset int) (uint16, error)

	// Int16 parses and returns an int16 value from the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.
	Int16(p []byte, offset int) (int16, error)

	// Uint32 parses and returns an uint32 value from the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.
	Uint32(p []byte, offset int) (uint32, error)

	// Int32 parses and returns an int32 value from the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.
	Int32(p []byte, offset int) (int32, error)

	// Float32 parses and returns a float32 value from the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.
	Float32(p []byte, offset int) (float32, error)

	// Uint64 parses and returns an uint64 value from the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.
	Uint64(p []byte, offset int) (uint64, error)

	// Int64 parses and returns an int64 value from the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.
	Int64(p []byte, offset int) (int64, error)

	// Float64 parses and returns a float64 value from the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.
	Float64(p []byte, offset int) (float64, error)

	// Close closes the underlying TCP connection. Returns a modbusclient.ErrNotconnected if the client is not connected to the server.
	Close() error
}

type client struct {
	Addr        string
	ConnTimeout time.Duration
	ByteOrder   binary.ByteOrder
	conn        net.Conn
}

// NewClient creates and returns a new Modbus TCP client.
func NewClient(addr string, connTimeout time.Duration, byteOrder binary.ByteOrder) Client {
	return &client{
		Addr:        addr,
		ConnTimeout: connTimeout,
		ByteOrder:   byteOrder,
	}
}

func (c *client) Connect() error {
	conn, err := net.DialTimeout("tcp4", c.Addr, c.ConnTimeout)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *client) SetDeadline(t time.Time) error {
	if c.conn == nil {
		return ErrNotConnected
	}

	return c.conn.SetDeadline(t)
}

func (c *client) Read(p []byte, unitID byte, addr uint16, count uint16) (int, error) {
	if c.conn == nil {
		return 0, ErrNotConnected
	}

	req := makeReadReq(unitID, addr, count)
	if _, err := c.conn.Write(req); err != nil {
		return 0, err
	}

	return c.conn.Read(p)
}

func (c *client) ReadErr(p []byte) error {
	if len(p) < readResHeaderLen {
		return ErrShortResponse
	}

	if errCode := p[errCodeIndex]; errCode != readFuncCode {
		excCode := p[excCodeIndex]
		return fmt.Errorf("%w, 0x%02x, 0x%02x", ErrModbusError, errCode, excCode)
	}

	return nil
}

func (c *client) Uint16(p []byte, offset int) (uint16, error) {
	offset += readResHeaderLen
	if len(p) < offset+2 {
		return 0, ErrShortPayload
	}

	return c.ByteOrder.Uint16(p[offset : offset+2]), nil
}

func (c *client) Int16(p []byte, offset int) (int16, error) {
	offset += readResHeaderLen
	if len(p) < offset+2 {
		return 0, ErrShortPayload
	}

	r := bytes.NewReader(p[offset : offset+2])
	var v int16
	if err := binary.Read(r, c.ByteOrder, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func (c *client) Uint32(p []byte, offset int) (uint32, error) {
	offset += readResHeaderLen
	if len(p) < offset+4 {
		return 0, ErrShortPayload
	}

	return c.ByteOrder.Uint32(p[offset : offset+4]), nil
}

func (c *client) Int32(p []byte, offset int) (int32, error) {
	offset += readResHeaderLen
	if len(p) < offset+4 {
		return 0, ErrShortPayload
	}

	r := bytes.NewReader(p[offset : offset+4])
	var v int32
	if err := binary.Read(r, c.ByteOrder, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func (c *client) Float32(p []byte, offset int) (float32, error) {
	offset += readResHeaderLen
	if len(p) < offset+4 {
		return 0, ErrShortPayload
	}

	r := bytes.NewReader(p[offset : offset+4])
	var v float32
	if err := binary.Read(r, c.ByteOrder, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func (c *client) Uint64(p []byte, offset int) (uint64, error) {
	offset += readResHeaderLen
	if len(p) < offset+8 {
		return 0, ErrShortPayload
	}

	return c.ByteOrder.Uint64(p[offset : offset+8]), nil
}

func (c *client) Int64(p []byte, offset int) (int64, error) {
	offset += readResHeaderLen
	if len(p) < offset+8 {
		return 0, ErrShortPayload
	}

	r := bytes.NewReader(p[offset : offset+8])
	var v int64
	if err := binary.Read(r, c.ByteOrder, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func (c *client) Float64(p []byte, offset int) (float64, error) {
	offset += readResHeaderLen
	if len(p) < offset+8 {
		return 0, ErrShortPayload
	}

	r := bytes.NewReader(p[offset : offset+8])
	var v float64
	if err := binary.Read(r, c.ByteOrder, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func (c *client) Close() error {
	if c.conn == nil {
		return ErrNotConnected
	}

	return c.conn.Close()
}

func makeReadReq(unitID byte, addr uint16, count uint16) []byte {
	return []byte{
		0x00,                      // transaction id, high
		0x00,                      // transcation id, low
		0x00,                      // protocol id, high
		0x00,                      // protocol id, low
		0x00,                      // length, high
		0x06,                      // length, low
		unitID,                    // unit id
		readFuncCode,              // function code
		byte((addr >> 8) & 0xFF),  // address, high
		byte(addr & 0xFF),         // address, low
		byte((count >> 8) & 0xFF), // register count, high
		byte(count & 0xFF),        // register count, low
	}
}
