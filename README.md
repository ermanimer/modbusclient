# modbusclient
modbusclient is a simple Modbus TCP client based on [Modbus Application Protocol Specification V1.1b3](https://www.modbus.org/docs/Modbus_Application_Protocol_V1_1b3.pdf)

# Supported Modbus Functions

- Read Holding Registers (0x03)

# Supported Data Types

- uint16
- int16
- uint32
- int32
- float32
- uint64
- int64
- float64

# Installation

```bash
go get -u github.com/ermanimer/modbusclient
```

# Methods

- **Connect() error:** Connect uses net.DialTimeout to establish an underlying TCP connection with the Modbus TCP server.

- **SetDeadline(t time.Time) error:** SetDeadline sets the underlying TCP connection's deadline. Returns a modbusclient.ErrNotconnected if the client is not connected.

- **Read(p []byte, unitID byte, addr uint16, count uint16) (n int, err error):** Read reads data from the Holding Registers of a Modbus TCP device and writes it to the provided payload. Returns the read-byte count and a modbusclient.ErrNotconnected if the client is not connected to the server.
	
- **ReadErr(p []byte) error:** ReadErr parses and returns the Modbus read error of the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.

- **Uint16(p []byte, offset int) (uint16, error):** Uint16 parses and returns an uint16 value from the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.

- **Int16(p []byte, offset int) (int16, error):** Int16 parses and returns an int16 value from the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.

- **Uint32(p []byte, offset int) (uint32, error):** Uint32 parses and returns an uint32 value from the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.

- **Int32(p []byte, offset int) (int32, error):** Int32 parses and returns an int32 value from the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.

- **Float32(p []byte, offset int) (float32, error):** Float32 parses and returns a float32 value from the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.

- **Uint64(p []byte, offset int) (uint64, error):** Uint64 parses and returns an uint64 value from the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.

- **Int64(p []byte, offset int) (int64, error):** nt64 parses and returns an int64 value from the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.

- **Float64(p []byte, offset int) (float64, error):** Float64 parses and returns a float64 value from the provided payload. Returns a modbusclient.ErrShortResponse if the payload is short.

- **Close() error:** Close closes the underlying TCP connection. Returns a modbusclient.ErrNotconnected if the client is not connected to the server.

# Sample Application

The sample application demonstrates reading a sample value from a Modbus device. 

```go
package main

import (
	"encoding/binary"
	"log"
	"time"

	"github.com/ermanimer/modbusclient"
)

// configurations
const (
	address          = "192.168.0.1:502" // address of the device
	connTimeout      = 5 * time.Second   // connection timeout
	unitID           = 0                 // unit id of the device
	startingAddresss = 0                 // starting address
	registerCount    = 2                 // register count
)

var byteOrder = binary.BigEndian // byte order of the Modbus TCP server

func main() {
	// create client
	client := modbusclient.NewClient(address, connTimeout, byteOrder)

	// connect
	err := client.Connect()
	if err != nil {
		log.Print(err)
		return
	}
	defer func() {
		err := client.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	// read
	buf := make([]byte, 256)
	n, err := client.Read(buf, unitID, startingAddresss, registerCount)
	if err != nil {
		log.Print(err)
		return
	}
	response := buf[:n]

	// check read error
	if err := client.ReadErr(response); err != nil {
		log.Print(err)
		return
	}

	// parse value
	value, err := client.Float32(response, 0)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("value: %.2f\n", value)
}
```
