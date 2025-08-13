package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// signinHandler обрабатывает аутентификацию пользователя.
func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{"error": "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		writeJSON(w, map[string]string{"error": "Пароль не настроен на сервере"}, http.StatusInternalServerError)
		return
	}

	var req struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]string{"error": "неверный формат запроса"}, http.StatusBadRequest)
		return
	}

	if req.Password != password {
		writeJSON(w, map[string]string{"error": "неверный пароль"}, http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"pass_hash": fmt.Sprintf("%x", sha256.Sum256([]byte(password))),
		"exp":       time.Now().Add(8 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(password))
	if err != nil {
		writeJSON(w, map[string]string{"error": "не удалось создать токен"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"token": tokenString}, http.StatusOK)
}

// auth — middleware для проверки аутентификации.
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		password := os.Getenv("TODO_PASSWORD")
		if password == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			writeJSON(w, map[string]string{"error": "Требуется аутентификация"}, http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неожиданный метод подписи: %v", t.Header["alg"])
			}
			return []byte(password), nil
		})

		if err != nil {
			writeJSON(w, map[string]string{"error": "Невалидный токен"}, http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			passHash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
			if claims["pass_hash"] != passHash {
				writeJSON(w, map[string]string{"error": "Невалидный токен"}, http.StatusUnauthorized)
				return
			}
		} else {
			writeJSON(w, map[string]string{"error": "Невалидный токен"}, http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
