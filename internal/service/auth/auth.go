package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dwnGnL/testWork/internal/config"
	"github.com/dwnGnL/testWork/lib/token"
)

type auth struct {
	conf *config.Config

	jwtClient token.JwtToken[*AdminAccessDetails]
}

type AdminAccessDetails struct {
	ID   int64  `json:"id"`
	User string `json:"user"`
	jwt.StandardClaims
}

func (a *AdminAccessDetails) Valid() error {
	a.VerifyExpiresAt(time.Now().Unix(), true)
	return nil
}

func (a *AdminAccessDetails) SetExpairesAt(exp int64) error {
	a.StandardClaims.ExpiresAt = exp
	return nil
}

func New(conf *config.Config) *auth {
	return &auth{
		conf:      conf,
		jwtClient: token.New[*AdminAccessDetails](conf.PrivKey, time.Duration(conf.ExpTokenSec*int64(time.Second))),
	}
}

func (s auth) Login(ctx context.Context, username string, password string) (string, error) {
	if username != "root" && password != "toor" {
		return "", errors.New("login/password incorect")
	}
	tokenStr, err := s.jwtClient.GenerateToken(&AdminAccessDetails{
		ID:   1,
		User: username,
	})
	if err != nil {
		return "", fmt.Errorf("generate token err %w", err)
	}
	return tokenStr, nil
}

func (s auth) CheckToken(tokenStr string) (int64, error) {
	claim, err := s.jwtClient.ExtractTokenMetadata(tokenStr)
	if err != nil {
		return 0, err
	}

	return claim.ID, nil
}
