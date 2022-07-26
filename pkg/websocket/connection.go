package websocket

import (
	"github.com/gobwas/ws"
	"net"
)

type Connection struct {
	net.Conn
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{conn}
}

func (conn *Connection) Send(code ws.OpCode, payload []byte) error {
	frame := ws.NewFrame(code, true, payload)
	return ws.WriteFrame(conn.Conn, frame)
}

func (conn *Connection) Read() (ws.Frame, error) {
	frame, err := ws.ReadFrame(conn.Conn)
	return frame, err
}

func WriteFrame(conn net.Conn, code ws.OpCode, payload []byte) error {
	frame := ws.NewFrame(code, true, payload)
	return ws.WriteFrame(conn, frame)
}

func ReadFrame(conn net.Conn) (ws.Frame, error) {
	return ws.ReadFrame(conn)
}
