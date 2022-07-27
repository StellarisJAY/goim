package service

import (
	context "context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/db/dao"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/snowflake"
	"github.com/stellarisJAY/goim/pkg/stringutil"
	"log"
	"time"
)

const (
	TokenTimeout = time.Hour * 2
)

var secretKey = []byte("secret key not defined yet")
var userIdGenerator = snowflake.NewSnowflake(config.Config.MachineID)

type Claims struct {
	jwt.StandardClaims
	UserId   string `json:"userId"`
	DeviceId string `json:"deviceId"`
}

// AuthServiceImpl 授权服务实现
type AuthServiceImpl struct {
}

func (as *AuthServiceImpl) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// 生成随机的Salt 和 密码的MD5
	salt := stringutil.RandomString(16)
	pwdMd5 := md5.Sum([]byte(request.Password + salt))
	err := dao.InsertUser(&model.User{
		ID:       userIdGenerator.NextID(),
		Account:  request.Account,
		Password: stringutil.HexString(pwdMd5[:]),
		NickName: request.NickName,
		Salt:     salt,
	})
	response := new(pb.RegisterResponse)
	if err != nil {
		response.Code = pb.Error
		response.Message = err.Error()
	} else {
		response.Code = pb.Success
	}
	return response, nil
}

// AuthorizeDevice 为用户设备授权，检查请求的密码，最终生成一个与用户ID和设备ID绑定的Token
func (as *AuthServiceImpl) AuthorizeDevice(ctx context.Context, request *pb.AuthRequest) (*pb.AuthResponse, error) {
	user, exists, err := dao.FindUserByAccount(request.Account)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &pb.AuthResponse{Code: pb.NotFound, Message: "user not found"}, nil
	}
	if !verifyPassword(user.Password, request.Password, user.Salt) {
		return &pb.AuthResponse{Code: pb.AccessDenied, Message: "wrong password"}, nil
	}
	// 保存授权记录
	_ = dao.InsertUserLoginLog(&model.DeviceLogin{
		UserID:    user.ID,
		DeviceID:  request.DeviceID,
		Timestamp: time.Now().UnixMilli(),
		Ip:        request.Ip,
	})
	// 生成 Token
	if token, err := generateToken(user.ID, request.DeviceID); err != nil {
		return nil, err
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
	// 检查Token是否过期
	if time.Now().After(time.UnixMilli(claims.ExpiresAt)) {
		log.Println("token expired: ", claims.ExpiresAt)
		return &pb.LoginResponse{
			Code: pb.AccessDenied,
		}, nil
	}
	if claims.DeviceId != request.DeviceID || claims.UserId != fmt.Sprintf("%x", request.UserID) {
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
func generateToken(userId int64, deviceId string) (string, error) {
	expireTime := time.Now().Add(TokenTimeout).UnixMilli()
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,
		},
		UserId:   fmt.Sprintf("%x", userId),
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

func verifyPassword(encoded, password, salt string) bool {
	sum := md5.Sum([]byte(password + salt))
	return stringutil.HexString(sum[:]) == encoded
}

// newSession 在 session管理器中记录 用户设备ID 与 网关服务器的绑定关系
func newSession(gateway, channel, userId, deviceId string) error {
	return nil
}
