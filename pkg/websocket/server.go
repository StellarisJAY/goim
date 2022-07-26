package websocket

import (
	"github.com/gobwas/ws"
	"log"
	"net"
	"net/http"
	"sync"
)

// Acceptor 用于验证websocket连接，具体在网关服务中验证连接的设备ID、用户ID、用户Token，最终返回一个channelID
type Acceptor interface {
	Accept(conn net.Conn, ctx AcceptorContext) (string, error)
}

type AcceptorContext struct {
	Gateway string
}

type Server struct {
	Address  string
	Channels sync.Map
	Acceptor Acceptor
}

func NewServer(address string) *Server {
	return &Server{
		Address:  address,
		Channels: sync.Map{},
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Printf("failed to upgrade HTTP to websocket for %s , error: %v", r.RemoteAddr, err)
			return
		}
		channelID, err := s.Acceptor.Accept(conn, AcceptorContext{Gateway: s.Address})
		if err != nil {
			_ = conn.Close()
			log.Printf("connection refused: %v", err)
			return
		}
		connection := NewConnection(conn)
		channel := NewChannel(connection)
		s.Channels.Store(channelID, channel)
		go func() {
			_ = channel.Start()
		}()
	})
	return http.ListenAndServe(s.Address, mux)
}

func (s *Server) Close() {
	s.Channels.Range(func(key, value interface{}) bool {
		value.(*Channel).Close()
		return true
	})
}
