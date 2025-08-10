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
	// проверяем, что используется метод POST
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// получаем пароль из переменной окружения
	password := os.Getenv("TODO_PASSWORD")
	if password == "" {
		http.Error(w, "Пароль не настроен на сервере", http.StatusInternalServerError)
		return
	}

	var req struct {
		Password string `json:"password"`
	}

	// декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]string{"error": "неверный формат запроса"})
		return
	}

	// сверяем пароли
	if req.Password != password {
		writeJSON(w, map[string]string{"error": "неверный пароль"})
		return
	}

	// создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"pass_hash": fmt.Sprintf("%x", sha256.Sum256([]byte(password))),
		"exp":       time.Now().Add(8 * time.Hour).Unix(),
	})

	// подписываем токен тем же паролем
	tokenString, err := token.SignedString([]byte(password))
	if err != nil {
		http.Error(w, "не удалось создать токен", http.StatusInternalServerError)
		return
	}

	// отправляем токен клиенту
	writeJSON(w, map[string]string{"token": tokenString})
}

// auth — middleware для проверки аутентификации.
func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// получаем пароль из переменной окружения
		password := os.Getenv("TODO_PASSWORD")
		if password == "" {
			// если пароль не установлен, пропускаем проверку
			next(w, r)
			return
		}

		// получаем токен из cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
			return
		}

		// парсим токен
		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
			// проверяем метод подписи
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неожиданный метод подписи: %v", t.Header["alg"])
			}
			return []byte(password), nil
		})

		if err != nil {
			http.Error(w, "Невалидный токен", http.StatusUnauthorized)
			return
		}

		// проверяем, что токен валиден и содержит нужные данные
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// сверяем хэш пароля из токена с хэшем текущего пароля
			passHash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
			if claims["pass_hash"] != passHash {
				http.Error(w, "Невалидный токен", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, "Невалидный токен", http.StatusUnauthorized)
			return
		}

		// если все проверки пройдены, передаем управление следующему обработчику
		next(w, r)
	}
}
