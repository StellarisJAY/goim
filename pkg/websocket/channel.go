package websocket

import (
	"context"
	"errors"
	"github.com/gobwas/ws"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"sync/atomic"
)

const (
	StatusNew uint32 = iota
	StatusStarted
	StatusClosed
)

type Channel struct {
	connection *Connection
	writeChan  chan []byte
	closed     context.Context
	cancel     context.CancelFunc
	status     uint32
	id         string
	userID     int64
	deviceID   string
}

func NewChannel(connection *Connection, id string, userID int64, deviceID string) *Channel {
	c := &Channel{
		connection: connection,
		writeChan:  make(chan []byte, 1<<20),
		status:     StatusNew,
		id:         id,
		userID:     userID,
		deviceID:   deviceID,
	}
	c.closed, c.cancel = context.WithCancel(context.Background())
	return c
}

func (c *Channel) Start() error {
	if !atomic.CompareAndSwapUint32(&c.status, StatusNew, StatusStarted) {
		return errors.New("can't start channel from current state")
	}
	go func() {
		err := c.writeLoop()
		if err != nil {
			log.Warn("write loop error: %v", err)
		}
		c.Close()
	}()
	<-c.closed.Done()
	c.gracefulShutdown()
	return nil
}

func (c *Channel) writeLoop() error {
	for payload := range c.writeChan {
		err := c.connection.Send(ws.OpBinary, payload)
		if err != nil {
			return err
		}
		n := len(c.writeChan)
		for i := 0; i < n; i++ {
			payload = <-c.writeChan
			err := c.connection.Send(ws.OpBinary, payload)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Channel) Push(payload []byte) error {
	if atomic.LoadUint32(&c.status) != StatusStarted {
		return errors.New("channel not available")
	}
	c.writeChan <- payload
	return nil
}

func (c *Channel) Close() {
	if !atomic.CompareAndSwapUint32(&c.status, StatusStarted, StatusClosed) {
		return
	}
	c.cancel()
}

func (c *Channel) gracefulShutdown() {
	// close channel and connection
	defer c.connection.Close()
	defer close(c.writeChan)
	// RPC call kick session
	conn, err := naming.GetClientConn("auth")
	if err != nil {
		return
	}
	client := pb.NewAuthClient(conn)
	resp, err := client.KickSession(context.TODO(), &pb.KickSessionRequest{
		UserID:   c.userID,
		DeviceID: c.deviceID,
	})
	if err != nil {
		log.Error(err)
		return
	}
	if resp.Code != pb.Success {
		log.Warn("graceful shutdown kick session error: %s", resp.Message)
		return
	}
}

func (c *Channel) Available() bool {
	return atomic.LoadUint32(&c.status) == StatusStarted
}

func (c *Channel) ID() string {
	return c.id
}
