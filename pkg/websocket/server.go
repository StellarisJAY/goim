package websocket

import (
	"github.com/gobwas/ws"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
	"net"
	"net/http"
	"sync"
)

// Acceptor 用于验证websocket连接，具体在网关服务中验证连接的设备ID、用户ID、用户Token，最终返回一个channelID
type Acceptor interface {
	Accept(conn net.Conn, ctx AcceptorContext) AcceptorResult
}

type AcceptorResult struct {
	UserID    int64
	DeviceID  string
	Error     error
	ChannelID string
}

type AcceptorContext struct {
	Gateway string
}

type Server struct {
	Address   string
	Channels  sync.Map
	Acceptor  Acceptor
	UserConns *sync.Map
}

func NewServer(address string) *Server {
	return &Server{
		Address:   address,
		Channels:  sync.Map{},
		UserConns: &sync.Map{},
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Warn("failed to upgrade HTTP to websocket for %s , error: %v", r.RemoteAddr, err)
			return
		}
		result := s.Acceptor.Accept(conn, AcceptorContext{Gateway: config.Config.RpcServer.Address})
		if result.Error != nil {
			log.Warn("connection from %s failed: %v", conn.RemoteAddr().String(), result.Error)
			_ = conn.Close()
			return
		}
		connection := NewConnection(conn)
		channel := NewChannel(connection, result.ChannelID, result.UserID, result.DeviceID)
		s.Channels.Store(channel.id, channel)
		s.UserConns.Store(result.UserID, channel)
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
