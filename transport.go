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
	// SendCommand encodes passed command and encodes response with error check.
	//
	// Response destination is passed to "out" and should be pointer.
	//
	// nil "out" value can be passed if command doesn't returns any response.
	SendCommand(conn net.Conn, cmd Command, out AbstractResponse) error
}

type JSONTransport struct{}

// NewJSONTransport returns JSON encoding/decoding transport
func NewJSONTransport() JSONTransport {
	return JSONTransport{}
}

// SendCommand implements Transport interface
func (t JSONTransport) SendCommand(conn net.Conn, cmd Command, out AbstractResponse) error {
	if err := json.NewEncoder(conn).Encode(cmd); err != nil {
		return err
	}

	result, err := bufio.NewReader(conn).ReadBytes(0x00)
	if err != nil && err != io.EOF {
		return err
	}

	trimmed := bytes.TrimRight(result, "\x00")
	// fix incorrect json response from miner ("}{")
	if cmd.Command == "stats" {
		trimmed = bytes.Replace(trimmed, []byte("}{"), []byte(","), 1)
	}

	isEmpty := out == nil
	if isEmpty {
		if len(trimmed) == 0 {
			return nil
		}
		out = new(GenericResponse)
	}

	if err := json.Unmarshal(trimmed, out); err != nil {
		if isEmpty {
			// just omit error if consumer passed empty response output
			return nil
		}
		return err
	}

	return out.HasError()
}
