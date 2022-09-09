package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/google/uuid"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/websocket"
	"google.golang.org/protobuf/proto"
	"net"
	"time"
)

const (
	HandshakeReadTimeout = time.Second * 5
)

// GateAcceptor 网关握手服务，负责接收websocket连接 并检查权限
type GateAcceptor struct {
}

func (acceptor *GateAcceptor) Accept(conn net.Conn, ctx websocket.AcceptorContext) websocket.AcceptorResult {
	_ = conn.SetReadDeadline(time.Now().Add(HandshakeReadTimeout))
	frame, err := websocket.ReadFrame(conn)
	if err != nil {
		return websocket.AcceptorResult{Error: fmt.Errorf("read handshake frame error: %w", err)}
	}
	if frame.Header.OpCode != ws.OpBinary {
		return websocket.AcceptorResult{Error: errors.New("wrong type of op code for handshake")}
	}
	if frame.Header.Masked {
		ws.Cipher(frame.Payload, frame.Header.Mask, 0)
		frame.Header.Masked = false
	}
	// 解码握手请求
	request := new(pb.HandshakeRequest)
	err = proto.Unmarshal(frame.Payload, request)
	if err != nil {
		return websocket.AcceptorResult{Error: errors.New("wrong binary content for handshake")}
	}
	// 生成ChannelID
	channel := generateChannelID()
	// 发送RPC请求，验证登录信息
	loginResp, err := login(request, ctx, channel)
	if err != nil {
		return websocket.AcceptorResult{Error: err}
	}
	response := new(pb.HandshakeResponse)
	response.Status = pb.HandshakeStatus_AccessDenied
	marshal, err := proto.Marshal(response)
	if err != nil {
		return websocket.AcceptorResult{Error: err}
	}
	frame = ws.NewFrame(ws.OpBinary, true, marshal)
	if err := ws.WriteFrame(conn, frame); err != nil {
		return websocket.AcceptorResult{Error: err}
	}
	log.Info("Accepted websocket connection, userID: %d, Device: %d, channel:  %s", loginResp.UserID, loginResp.DeviceID, channel)
	return websocket.AcceptorResult{
		UserID:    loginResp.UserID,
		DeviceID:  loginResp.DeviceID,
		ChannelID: channel,
	}

}

// login 设备登录，设备必须通过websocket发送握手包接入聊天服务，握手时会从授权服务检查用户Token和设备信息
func login(request *pb.HandshakeRequest, ctx websocket.AcceptorContext, channel string) (*pb.LoginResponse, error) {
	// 连接到授权服务
	conn, err := naming.GetClientConn("auth")
	if err != nil {
		return nil, err
	}
	client := pb.NewAuthClient(conn)
	// RPC 调用进行登录
	response, err := client.LoginDevice(context.Background(), &pb.LoginRequest{
		Token:   request.Token,
		Gateway: ctx.Gateway,
		Channel: channel,
	})
	if err != nil {
		return nil, err
	}
	switch response.Code {
	case pb.Success:
		return response, nil
	case pb.Error:
		log.Errorf("login error: %s", response.Message)
		return nil, errors.New("access denied")
	default:
		return nil, errors.New("access denied")
	}
}

// generateChannelID 暂时使用 UUID 作为channelID
func generateChannelID() string {
	uid, _ := uuid.NewUUID()
	return uid.String()
}
