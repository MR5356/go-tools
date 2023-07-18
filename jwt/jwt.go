package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	InvalidTokenError = errors.New("invalid token")
)

type Service[T any] struct {
	signKey        []byte
	issuer         string
	model          T
	expireDuration time.Duration
}

type CustomClaims[T any] struct {
	Model T
	jwt.RegisteredClaims
}

// NewJWTService 创建一个JWTService
func NewJWTService[T any](secret, issuer string, expireDuration time.Duration, model T) *Service[T] {
	return &Service[T]{
		signKey:        []byte(secret),
		issuer:         issuer,
		model:          model,
		expireDuration: expireDuration,
	}
}

// CreateToken 新建一个Token
func (s *Service[T]) CreateToken(model T) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		CustomClaims[T]{
			Model: model,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(s.expireDuration)),
				NotBefore: jwt.NewNumericDate(now.Add(-100 * time.Second)),
				Issuer:    s.issuer,
			},
		},
	)
	tokenString, err := token.SignedString(s.signKey)
	return tokenString, err
}

// ParseTokenIgnoreExpired 解析Token，返回结构体，忽略Token过期错误
func (s *Service[T]) ParseTokenIgnoreExpired(tokenString string) (model T, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims[T]{}, func(token *jwt.Token) (interface{}, error) {
		return s.signKey, nil
	})
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return
	}

	claims, ok := token.Claims.(*CustomClaims[T])
	if !ok || claims.Issuer != s.issuer {
		err = InvalidTokenError
		return
	}
	return claims.Model, nil
}

// ParseToken 解析Token，返回结构体
func (s *Service[T]) ParseToken(tokenString string) (model T, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims[T]{}, func(token *jwt.Token) (interface{}, error) {
		return s.signKey, nil
	})
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*CustomClaims[T])
	if !ok || !token.Valid || claims.Issuer != s.issuer {
		err = InvalidTokenError
		return
	}

	return claims.Model, nil
}

// CreateTokenWithRefreshToken 新建一个Token，同时带有刷新Token用的refreshToken
func (s *Service[T]) CreateTokenWithRefreshToken(model T, refreshTokenExpireDuration time.Duration) (tokenString, refreshTokenString string, err error) {
	tokenString, err = s.CreateToken(model)
	if err != nil {
		return
	}
	refreshTokenString, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpireDuration)),
		Issuer:    s.issuer,
	}).SignedString(s.signKey)
	return
}

// RefreshToken 使用refreshToken刷新Token
func (s *Service[T]) RefreshToken(tString, rString string, refreshTokenExpireDuration time.Duration) (tokenString, refreshTokenString string, err error) {
	if _, err = jwt.Parse(rString, func(token *jwt.Token) (interface{}, error) {
		return s.signKey, nil
	}); err != nil {
		return
	}

	model, err := s.ParseTokenIgnoreExpired(tString)
	if err != nil {
		return
	}

	return s.CreateTokenWithRefreshToken(model, refreshTokenExpireDuration)
}
