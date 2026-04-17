package teamserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/z3vxo/kronos/internal/config"
)

func authMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")

		if authToken == "" {
			SendJSONError(w, "missing token", http.StatusUnauthorized)

			return
		}

		tokenStr := strings.TrimPrefix(authToken, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("incorrect siging method")
			}
			return []byte(config.Cfg.TS.Auth.JwtSecret), nil
		})

		if err != nil || !token.Valid {
			SendJSONError(w, "invalid token", http.StatusUnauthorized)

			return
		}
		next.ServeHTTP(w, r)
	})
}

func CheckLogin(user, pass string) bool {
	return user == config.Cfg.TS.Auth.Username && pass == config.Cfg.TS.Auth.Password

}

func CraftJWT(user string) (string, error) {
	claims := jwt.MapClaims{
		"user": user,
		"exp":  config.Cfg.TS.Auth.TokenHours,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.Cfg.TS.Auth.JwtSecret))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return tokenString, nil
}
