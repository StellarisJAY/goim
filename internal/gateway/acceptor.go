package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/log"
	"github.com/stellarisJAY/goim/pkg/naming"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/websocket"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"net"
	"time"
)

const (
	HandshakeReadTimeout = time.Second * 5
)

// GateAcceptor 网关握手服务，负责接收websocket连接 并检查权限
type GateAcceptor struct {
	globalTracer opentracing.Tracer
}

func (acceptor *GateAcceptor) Accept(conn net.Conn, ctx websocket.AcceptorContext) websocket.AcceptorResult {
	_ = conn.SetReadDeadline(time.Now().Add(HandshakeReadTimeout))
	frame, err := websocket.ReadFrame(conn)
	if err != nil {
		return websocket.AcceptorResult{Error: fmt.Errorf("read handshake frame error: %w", err)}
	}
	if frame.Header.Masked {
		ws.Cipher(frame.Payload, frame.Header.Mask, 0)
		frame.Header.Masked = false
	}
	// 解码握手请求
	request, err := unmarshalHandshakeRequest(frame.Payload)
	if err != nil {
		return websocket.AcceptorResult{Error: fmt.Errorf("unmarshal handshake request error: %w", err)}
	}
	// 生成ChannelID
	channel := generateChannelID()
	// 发送RPC请求，验证登录信息
	loginResp, err := acceptor.login(request, ctx, channel)
	if err != nil {
		return websocket.AcceptorResult{Error: err}
	}
	response := new(pb.HandshakeResponse)
	response.Status = pb.HandshakeStatus_AccessDenied
	marshal, err := marshalHandshakeResponse(response)
	if err != nil {
		return websocket.AcceptorResult{Error: err}
	}
	frame = ws.NewFrame(ws.OpText, true, marshal)
	if err := ws.WriteFrame(conn, frame); err != nil {
		return websocket.AcceptorResult{Error: err}
	}
	log.Info("Accepted websocket connection",
		zap.Int64("userID", loginResp.UserID),
		zap.String("deviceID", loginResp.DeviceID),
		zap.String("channel", channel))
	return websocket.AcceptorResult{
		UserID:    loginResp.UserID,
		DeviceID:  loginResp.DeviceID,
		ChannelID: channel,
	}

}

// login 设备登录，设备必须通过websocket发送握手包接入聊天服务，握手时会从授权服务检查用户Token和设备信息
func (acceptor *GateAcceptor) login(request *pb.HandshakeRequest, ctx websocket.AcceptorContext, channel string) (*pb.LoginResponse, error) {
	// 连接到授权服务
	conn, err := naming.GetClientConn(pb.UserServiceName)
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
		log.Warn("login request failed", zap.Error(err))
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

func unmarshalHandshakeRequest(payload []byte) (*pb.HandshakeRequest, error) {
	request := &pb.HandshakeRequest{}
	if config.Config.Gateway.UseJsonMsg {
		if err := json.Unmarshal(payload, request); err != nil {
			return nil, err
		}
	} else {
		if err := proto.Unmarshal(payload, request); err != nil {
			return nil, err
		}
	}
	return request, nil
}

func marshalHandshakeResponse(response *pb.HandshakeResponse) ([]byte, error) {
	if config.Config.Gateway.UseJsonMsg {
		return json.Marshal(response)
	} else {
		return proto.Marshal(response)
	}
}
