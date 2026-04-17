package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/z3vxo/kronos/internal/config"
	"github.com/z3vxo/kronos/internal/httputil"
)

type Auth struct {
	User       string
	Paswd      string
	JwtSecret  []byte
	TokenHours int
}

func NewAuth(user, passwd, secret string, Hours int) *Auth {
	return &Auth{
		User:       user,
		Paswd:      passwd,
		JwtSecret:  []byte(secret),
		TokenHours: Hours,
	}
}

func (a *Auth) ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	if tokenStr == "" {
		return nil, errors.New("Missing token")
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Incorrect signing algo")
		}
		return a.JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("Invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}

func (a *Auth) AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := a.ValidateToken(r.Header.Get("Authorization"))
		if err != nil {
			httputil.SendJSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *Auth) CheckLogin(user, pass string) bool {
	return user == config.Cfg.TS.Auth.Username && pass == config.Cfg.TS.Auth.Password

}

func (a *Auth) CraftJWT(user string) (string, error) {
	claims := jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(time.Duration(a.TokenHours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.JwtSecret))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return tokenString, nil
}
