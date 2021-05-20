package util

import (
	"fmt"
	"github.com/bearname/videohost/videoserver/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"net/http"
	"os"
	"time"
)

func CreateTokenWithDuration(userKey string, username string, role model.Role, duration time.Duration) (string, error) {
	var err error
	claims := jwt.MapClaims{}
	claims["userId"] = userKey
	claims["role"] = role.Values()
	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(duration).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := getSecret()
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func CreateToken(userKey string, username string, role model.Role) (string, error) {
	return CreateTokenWithDuration(userKey, username, role, time.Minute*2)
}

func CheckToken(tokenString string) (*jwt.Token, bool) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(getSecret()), nil
	}
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, false
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, false
	}
	return token, true
}

func IsUsernameContextOk(username string, r *http.Request) bool {
	usernameCtx, ok := context.Get(r, "username").(string)
	if !ok {
		return false
	}
	if usernameCtx != username {
		return false
	}
	return true
}

func getSecret() string {
	secret := os.Getenv("ACCESS_SECRET")

	if secret == "" {
		secret = "sdmalncnjsdsmf"
	}

	return secret
}
