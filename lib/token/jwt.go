package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dwnGnL/testWork/lib/goerrors"
)

var (
	ERR_TOKEN_NOT_FOUND = errors.New("token not found")
	CLAIMS_NOT_VALID    = errors.New("claims not valid")
)

type JwtToken[claim MyClaim] struct {
	secretKey          string
	expirationDuration time.Duration
}

func New[claim MyClaim](key string, duration time.Duration) JwtToken[claim] {
	return JwtToken[claim]{secretKey: key, expirationDuration: duration}
}

type MyClaim interface {
	Valid() error
	SetExpairesAt(exp int64) error
}

func (j *JwtToken[claim]) verifyToken(tokenStr string) (claim, error) {
	var nilClaim claim
	goerrors.Log().Println("start verifyToken")

	goerrors.Log().Println("start jwt.Parse")
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nilClaim, err
	}
	goerrors.Log().Println("start token.Claims")

	if v, ok := token.Claims.(jwt.MapClaims); ok {
		jsonbody, err := json.Marshal(v)
		if err != nil {
			return nilClaim, err
		}

		if err := json.Unmarshal(jsonbody, &nilClaim); err != nil {
			return nilClaim, err
		}
		return nilClaim, nil
	} else {
		return nilClaim, CLAIMS_NOT_VALID
	}
	//return token, nil
}

func (j *JwtToken[claim]) ExtractTokenMetadata(tokenStr string) (claim, error) {
	var nilClaim claim
	goerrors.Log().Println("start ExtractTokenMetadata")
	tokenClaim, err := j.verifyToken(tokenStr)
	if err != nil {
		return nilClaim, err
	}
	goerrors.Log().Println("start Valid")

	if tokenClaim.Valid() != nil {
		goerrors.Log().Info("token not valid")
		return nilClaim, errors.New("token not valid")
	}
	goerrors.Log().Info("token: ", tokenClaim)

	return tokenClaim, err
}

func (j *JwtToken[claim]) GenerateToken(myClaim claim) (string, error) {
	err := myClaim.SetExpairesAt(
		time.Now().Add(j.expirationDuration).Unix(),
	)
	if err != nil {
		return "", fmt.Errorf("err on set standart data")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaim)
	return token.SignedString([]byte(j.secretKey))
}
