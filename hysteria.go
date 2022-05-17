// Package pt_hysteria provides a PT 3.0 Go API wrapper around the connections used by hysteria
package pt_hysteria

import (
	"context"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/tobyxdd/hysteria/pkg/transport"
)

// Create outgoing transport connection
// The Dial method implements the **Client Factory **abstract interface.

type HysteriaClient struct {
	serverAddress string
}

func NewHysteriaClient(serverAddress string) *HysteriaClient {
	return &HysteriaClient{serverAddress: serverAddress}
}

// func (client *HysteriaClient) Dial() (net.Conn, error) {
// 	return net.Dial(serverAddress)
// }

// Create listener for incoming transport connection
// The Listen method implements the **Server Factory **abstract interface.

type HysteriaServer struct {
	listenAddress string
}

func NewHysteriaServer(address string, obfs string, cert string, key string, mode string, password string) *HysteriaListener {
	hl, err := transport.DefaultServerTransport.QUICListen()

	return &HysteriaListener{hl}
}

// func (server *HysteriaServer) Listen() (net.Conn, error) {
// 	return NewHysteriaListener(listenAddress), nil
// }

type HysteriaListener struct {
	ql quic.Listener
}

type StreamConn struct {
	c quic.Connection
	s quic.Stream
}

func (hl *HysteriaListener) Listen() (net.Listener, error) {
	return hl, nil
}

func (hl *HysteriaListener) Accept() (sc net.Conn, e error) {
	for {
		qconn, e := hl.ql.Accept(context.Background())

		if e != nil {
			return nil, e
		}

		go func() {
			s, err := qconn.AcceptStream(context.Background())
			if err != nil {
				return
			}

			defer s.Close()

			sc = StreamConn{qconn, s}
			return
		}()
	}
}

// Addr returns the listener's network address.
func (hl *HysteriaListener) Addr() net.Addr {
	return hl.ql.Addr()
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (hl *HysteriaListener) Close() error {
	return hl.ql.Close()
}

// LocalAddr returns the local network address.
// needed to fulfill the net.Conn interface
func (sc StreamConn) LocalAddr() net.Addr {
	return sc.c.LocalAddr()
}

func (sc StreamConn) Close() error {
	return sc.s.Close()
}

func (sc StreamConn) Read(b []byte) (n int, err error) {
	return sc.s.Read(b)
}

// RemoteAddr returns the remote network address.
func (sc StreamConn) RemoteAddr() net.Addr {
	return sc.c.RemoteAddr()
}

func (sc StreamConn) SetDeadline(t time.Time) error {
	return sc.s.SetDeadline(t)
}

func (sc StreamConn) SetReadDeadline(t time.Time) error {
	return sc.s.SetReadDeadline(t)
}

func (sc StreamConn) SetWriteDeadline(t time.Time) error {
	return sc.s.SetWriteDeadline(t)
}

func (sc StreamConn) Write(b []byte) (n int, err error) {
	return sc.s.Write(b)
}
