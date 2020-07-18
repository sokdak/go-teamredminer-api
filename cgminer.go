package cgminer

import (
	"context"
	"fmt"
	"net"
	"time"
)

// ConnectError represents API connection error
type ConnectError struct {
	err error
}

func (err ConnectError) Error() string {
	return fmt.Sprintf("connect error (%s)", err.err)
}

func (err ConnectError) Unwrap() error {
	return err.err
}

// CGMiner is cgminer API client
type CGMiner struct {
	// Address is API endpoint address (host:port)
	Address string

	// Timeout is request timeout
	Timeout time.Duration

	// Dialer is network dialer
	Dialer net.Dialer

	// Transport is request and response decoder.
	//
	// CGMiner might have one of two API formats - JSON or plain text.
	// JSON is default one.
	Transport Transport
}

// Call sends command to cgminer API and writes result to passed response output
// or returns error.
//
// If command doesn't returns any response, nil "out" value should be passed.
//
// For context-based requests, use `CallContext()`
func (c *CGMiner) Call(cmd Command, out AbstractResponse) error {
	return c.CallContext(context.Background(), cmd, out)
}

// CallContext sends command to cgminer API using the provided context.
//
// If command doesn't returns any response, nil "out" value should be passed.
func (c *CGMiner) CallContext(ctx context.Context, cmd Command, out AbstractResponse) error {
	conn, err := c.Dialer.DialContext(ctx, "tcp", c.Address)
	if err != nil {
		return ConnectError{err: err}
	}

	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(c.Timeout))
	return c.Transport.SendCommand(conn, cmd, out)
}

// NewCGMiner returns a CGMiner client with JSON API transport
func NewCGMiner(hostname string, port int, timeout time.Duration) *CGMiner {
	return &CGMiner{
		Address:   fmt.Sprintf("%s:%d", hostname, port),
		Timeout:   timeout,
		Transport: NewJSONTransport(),
		Dialer: net.Dialer{
			Timeout: timeout,
		},
	}
}
