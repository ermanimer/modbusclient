package modbusclient

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestReadErr(t *testing.T) {
	c := &client{}

	var errCode byte = 0x83
	var excCode byte = 0x01
	p := []byte{
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x06,
		0x01,
		errCode,
		excCode,
	}
	err := c.ReadErr(p)
	if !errors.Is(err, ErrModbusError) {
		t.Error("read error is not ErrModbusError")
	}
	if !strings.Contains(err.Error(), fmt.Sprintf("0x%02x", errCode)) {
		t.Error("read error does not contain error code")
	}
	if !strings.Contains(err.Error(), fmt.Sprintf("0x%02x", excCode)) {
		t.Error("read error does not contain exception code")
	}
}

func TestErrShortPayload(t *testing.T) {
	c := &client{
		ByteOrder: binary.BigEndian,
	}

	var p []byte

	_, err := c.Uint16(p, 0)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}

	_, err = c.Int16(p, 0)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}

	_, err = c.Uint32(p, 0)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}

	_, err = c.Int32(p, 0)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}

	_, err = c.Float32(p, 0)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}

	_, err = c.Uint64(p, 0)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}

	_, err = c.Int64(p, 0)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}

	_, err = c.Float64(p, 0)
	if !errors.Is(err, ErrShortPayload) {
		t.Error("error is not ErrShortPayload")
	}
}

func TestUint16(t *testing.T) {
	c := &client{
		ByteOrder: binary.BigEndian,
	}

	var expected uint16 = 1
	p := make([]byte, readResHeaderLen+2)
	c.ByteOrder.PutUint16(p[readResHeaderLen:], expected)

	v, err := c.Uint16(p, 0)
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}

func TestInt16(t *testing.T) {
	c := &client{
		ByteOrder: binary.BigEndian,
	}

	var expected int16 = 1
	p := make([]byte, readResHeaderLen)
	w := bytes.NewBuffer(p)
	binary.Write(w, c.ByteOrder, expected)
	p = w.Bytes()

	v, err := c.Int16(p, 0)
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}

func TestUint32(t *testing.T) {
	c := &client{
		ByteOrder: binary.BigEndian,
	}

	var expected uint32 = 1
	p := make([]byte, readResHeaderLen+4)
	c.ByteOrder.PutUint32(p[readResHeaderLen:], expected)

	v, err := c.Uint32(p, 0)
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}

func TestInt32(t *testing.T) {
	c := &client{
		ByteOrder: binary.BigEndian,
	}

	var expected int32 = 1
	p := make([]byte, readResHeaderLen)
	w := bytes.NewBuffer(p)
	binary.Write(w, c.ByteOrder, expected)
	p = w.Bytes()

	v, err := c.Int32(p, 0)
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}

func TestFloat32(t *testing.T) {
	c := &client{
		ByteOrder: binary.BigEndian,
	}

	var expected float32 = 1
	p := make([]byte, readResHeaderLen)
	w := bytes.NewBuffer(p)
	binary.Write(w, c.ByteOrder, expected)
	p = w.Bytes()

	v, err := c.Float32(p, 0)
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}

func TestUint64(t *testing.T) {
	c := &client{
		ByteOrder: binary.BigEndian,
	}

	var expected uint64 = 1
	p := make([]byte, readResHeaderLen+8)
	c.ByteOrder.PutUint64(p[readResHeaderLen:], expected)

	v, err := c.Uint64(p, 0)
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}

func TestInt64(t *testing.T) {
	c := &client{
		ByteOrder: binary.BigEndian,
	}

	var expected int64 = 1
	p := make([]byte, readResHeaderLen)
	w := bytes.NewBuffer(p)
	binary.Write(w, c.ByteOrder, expected)
	p = w.Bytes()

	v, err := c.Int64(p, 0)
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}

func TestFloat64(t *testing.T) {
	c := &client{
		ByteOrder: binary.BigEndian,
	}

	var expected float64 = 1
	p := make([]byte, readResHeaderLen)
	w := bytes.NewBuffer(p)
	binary.Write(w, c.ByteOrder, expected)
	p = w.Bytes()

	v, err := c.Float64(p, 0)
	if err != nil {
		t.Error(err)
	}
	if v != expected {
		t.Error("value is not equal to expected", v, expected)
	}
}
