package service

import (
	context "context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/stellarisJAY/goim/pkg/db"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"log"
	"time"
)

const (
	TokenTimeout = time.Hour * 2
)

var secretKey = []byte("secret key not defined yet")

type Claims struct {
	jwt.StandardClaims
	UserId   string `json:"userId"`
	DeviceId string `json:"deviceId"`
}

// AuthServiceImpl 授权服务实现
type AuthServiceImpl struct {
}

func (as *AuthServiceImpl) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	tx := db.DB.MySQL.Create(&model.User{
		Account:   request.Account,
		Password:  request.Password,
		NickName:  request.NickName,
		CreatedAt: time.Now(),
	})
	response := new(pb.RegisterResponse)
	if tx.Error != nil {
		response.Code = pb.Error
		response.Message = tx.Error.Error()
	} else {
		response.Code = pb.Success
	}
	return response, nil
}

// AuthorizeDevice 为用户设备授权，检查请求的密码，最终生成一个与用户ID和设备ID绑定的Token
func (as *AuthServiceImpl) AuthorizeDevice(ctx context.Context, request *pb.AuthRequest) (*pb.AuthResponse, error) {
	// 检查用户ID 和 密码
	if !verifyPassword(request.UserId, request.Password) {
		return &pb.AuthResponse{
			Code: pb.WrongPassword,
		}, nil
	}
	log.Println("authorize: ", request)
	// 生成 Token
	if token, err := generateToken(request.UserId, request.DeviceId); err != nil {
		return &pb.AuthResponse{
			Code:  pb.Error,
			Token: "",
		}, nil
	} else {
		return &pb.AuthResponse{
			Code:  pb.Success,
			Token: token,
		}, nil
	}
}

// LoginDevice 设备登录，该方法在设备接入聊天服务前调用，用于检查Token是否合法，并最终记录session
func (as *AuthServiceImpl) LoginDevice(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 解析Token
	claims, err := parseToken(request.Token)
	if err != nil {
		log.Println(err)
		return &pb.LoginResponse{
			Code: pb.AccessDenied,
		}, nil
	}
	log.Println(claims)
	// 检查Token是否过期
	if time.Now().After(time.UnixMilli(claims.ExpiresAt)) {
		log.Println("token expired: ", claims.ExpiresAt)
		return &pb.LoginResponse{
			Code: pb.AccessDenied,
		}, nil
	}
	if claims.DeviceId != request.DeviceId || claims.UserId != request.UserId {
		log.Println(claims.DeviceId, ", ", claims.UserId)
		return &pb.LoginResponse{
			Code: pb.AccessDenied,
		}, nil
	}
	// 缓存设备的session信息，在 Redis中记录 userId: {deviceId: {gateway, channel}}
	if err := newSession(request.Gateway, request.Channel, claims.UserId, claims.DeviceId); err != nil {
		return &pb.LoginResponse{
			Code:    pb.Error,
			Message: err.Error(),
		}, nil
	}
	return &pb.LoginResponse{
		Code: pb.Success,
	}, nil
}

// generateToken 通过用户ID 和 设备ID 生成Token
func generateToken(userId, deviceId string) (string, error) {
	expireTime := time.Now().Add(TokenTimeout).UnixMilli()
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,
		},
		UserId:   userId,
		DeviceId: deviceId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// parseToken 解析Token
func parseToken(signed string) (*Claims, error) {
	claims := new(Claims)
	c, err := jwt.ParseWithClaims(signed, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if c.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func verifyPassword(userId, password string) bool {
	return true
}

// newSession 在 session管理器中记录 用户设备ID 与 网关服务器的绑定关系
func newSession(gateway, channel, userId, deviceId string) error {
	return nil
}
