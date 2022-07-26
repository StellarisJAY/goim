package gateway

import (
	"context"
	"errors"
	"github.com/gobwas/ws"
	"github.com/google/uuid"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/websocket"
	"google.golang.org/protobuf/proto"
	"net"
	"time"
)

const (
	HandshakeTimeout     = time.Second * 5
	HandshakeReadTimeout = time.Second * 5
)

// GateAcceptor 网关握手服务，负责接收websocket连接 并检查权限
type GateAcceptor struct {
	nextId string
}

func (acceptor *GateAcceptor) Accept(conn net.Conn, ctx websocket.AcceptorContext) (string, error) {
	_ = conn.SetReadDeadline(time.Now().Add(HandshakeReadTimeout))
	frame, err := websocket.ReadFrame(conn)
	if err != nil {
		return "", err
	}
	if frame.Header.OpCode != ws.OpBinary {
		return "", errors.New("wrong type of op code for handshake")
	}
	if frame.Header.Masked {
		ws.Cipher(frame.Payload, frame.Header.Mask, 0)
		frame.Header.Masked = false
	}
	// 解码握手请求
	request := new(pb.HandshakeRequest)
	err = proto.Unmarshal(frame.Payload, request)
	if err != nil {
		return "", errors.New("wrong binary content for handshake")
	}
	// 生成ChannelID
	channel := generateChannelID()
	// 发送RPC请求，验证登录信息
	if err := login(request, ctx, channel); err != nil {
		return "", err
	} else {
		return channel, nil
	}
}

// login 设备登录，设备必须通过websocket发送握手包接入聊天服务，握手时会从授权服务检查用户Token和设备信息
func login(request *pb.HandshakeRequest, ctx websocket.AcceptorContext, channel string) error {
	// 连接到授权服务
	conn, err := naming.GetClientConn("auth")
	if err != nil {
		return err
	}
	client := pb.NewAuthClient(conn)
	// RPC 调用进行登录
	response, err := client.LoginDevice(context.Background(), &pb.LoginRequest{
		Token:   request.Token,
		Gateway: ctx.Gateway,
		Channel: channel,
	})
	if err != nil {
		return err
	}
	switch response.Code {
	case pb.Success:
		return nil
	default:
		return errors.New("access denied")
	}
}

// generateChannelID 暂时使用 UUID 作为channelID
func generateChannelID() string {
	uid, _ := uuid.NewUUID()
	return uid.String()
}
