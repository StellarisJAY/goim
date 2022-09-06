package authutil

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/stellarisJAY/goim/pkg/config"
	"time"
)

type Claims struct {
	jwt.StandardClaims
	UserId   string `json:"userId"`
	DeviceId string `json:"deviceId"`
}

func ValidateToken(token string) (userID, deviceID string, valid, expired bool, expireAt int64) {
	claims, err := parseToken(token)
	if err != nil {
		valid = false
		return
	}
	if claims.ExpiresAt <= time.Now().UnixMilli() {
		expired = true
		return
	}
	expireAt = claims.ExpiresAt
	userID = claims.UserId
	deviceID = claims.DeviceId
	valid = true
	return
}

// parseToken 解析Token
func parseToken(signed string) (*Claims, error) {
	claims := new(Claims)
	c, err := jwt.ParseWithClaims(signed, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.TokenSecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if c.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
