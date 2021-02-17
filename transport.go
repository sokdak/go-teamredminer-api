package cgminer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net"
)

// AbstractResponse is generic command response which provides execution status
type AbstractResponse interface {
	// HasError returns status error
	HasError() error
}

// Transport encodes request and decodes response from cgminer API.
//
// CGMiner have two communication formats: JSON and plain-text.
// One of them or both might be enabled.
//
// Use corresponding transport for specific case.
type Transport interface {
	// SendCommand encodes and sends passed command
	SendCommand(conn net.Conn, cmd Command) error

	// DecodeResponse reads and decodes response from CGMiner API connection.
	//
	// Response destination is passed to "out" and should be pointer.
	//
	// nil "out" value can be passed if command doesn't returns any response.
	DecodeResponse(conn net.Conn, cmd Command, out AbstractResponse) error
}

var _ Transport = (*JSONTransport)(nil)

type JSONTransport struct{}

// NewJSONTransport returns JSON encoding/decoding transport
func NewJSONTransport() JSONTransport {
	return JSONTransport{}
}

// SendCommand implements Transport interface
func (t JSONTransport) SendCommand(conn net.Conn, cmd Command) error {
	return json.NewEncoder(conn).Encode(cmd)
}

// DecodeResponse implements Transport interface
func (t JSONTransport) DecodeResponse(conn net.Conn, cmd Command, out AbstractResponse) error {
	rsp, err := readWithNullTerminator(conn)
	if err != nil && err != io.EOF {
		return err
	}

	// fix incorrect json response from miner ("}{")
	if cmd.Command == "stats" {
		rsp = bytes.Replace(rsp, []byte("}{"), []byte(","), 1)
	}

	isEmpty := out == nil
	if isEmpty {
		if len(rsp) == 0 {
			return nil
		}
		out = new(GenericResponse)
	}

	if err := json.Unmarshal(rsp, out); err != nil {
		if isEmpty {
			// just omit error if consumer passed empty response output
			return nil
		}
		return err
	}

	return out.HasError()
}

// readWithNullTerminator reads cgminer response, but stops
// at null terminator (0x00)
func readWithNullTerminator(r io.Reader) ([]byte, error) {
	result, err := bufio.NewReader(r).ReadBytes(0x00)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return bytes.TrimRight(result, "\x00"), nil
}
