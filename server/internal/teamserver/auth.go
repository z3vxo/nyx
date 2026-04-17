package teamserver

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/z3vxo/kronos/internal/config"
)

func Send404(w http.ResponseWriter) {
	w.Header().Set("Server", "kronos")
	w.Header().Set("Content-Type", "text/html")

	path := config.Cfg.Server.NotFoundFile
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		path = home + path[1:]
	}

	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "<h1>404 not found</h1>", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write(content)
}

func authMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")

		if authToken == "" {
			//Send404(w)
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
			//Send404(w)
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
		"exp":  time.Now().Add(time.Duration(config.Cfg.TS.Auth.TokenHours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.Cfg.TS.Auth.JwtSecret))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return tokenString, nil
}
