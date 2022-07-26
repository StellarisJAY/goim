package websocket

import (
	"context"
	"errors"
	"github.com/gobwas/ws"
	"log"
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
}

func NewChannel(connection *Connection, id string) *Channel {
	c := &Channel{
		connection: connection,
		writeChan:  make(chan []byte, 1<<20),
		status:     StatusNew,
	}
	c.closed, c.cancel = context.WithCancel(context.Background())
	return c
}

func (c *Channel) Start() error {
	go func() {
		err := c.ReadLoop()
		if err != nil {
			log.Println("read loop error: ", err)
		}
		c.Close()
	}()
	go func() {
		err := c.writeLoop()
		if err != nil {
			log.Println("write loop error: ", err)
		}
		c.Close()
	}()
	<-c.closed.Done()
	c.gracefulShutdown()
	return nil
}

func (c *Channel) ReadLoop() error {
	for {
		select {
		case <-c.closed.Done():
			return nil
		default:
		}
		frame, err := c.connection.Read()
		if err != nil {
			return err
		}
		if frame.Header.OpCode == ws.OpPing {
			err = c.connection.Send(ws.OpPong, nil)
			if err != nil {
				return err
			}
		}
		if frame.Header.Masked {
			ws.Cipher(frame.Payload, frame.Header.Mask, 0)
			frame.Header.Masked = false
		}
		log.Printf("received frame: %v", frame)
	}
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
	_ = c.connection.Close()
	close(c.writeChan)
}

func (c *Channel) Available() bool {
	return atomic.LoadUint32(&c.status) == StatusStarted
}

func (c *Channel) ID() string {
	return c.id
}
