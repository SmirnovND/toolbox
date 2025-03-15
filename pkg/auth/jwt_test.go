package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Тест для getToken
func TestGetToken(t *testing.T) {
	tests := []struct {
		name      string
		setupReq  func() *http.Request
		wantToken string
		wantErr   bool
	}{
		{
			name: "Token from Authorization header",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "Bearer test-token")
				return req
			},
			wantToken: "test-token",
			wantErr:   false,
		},
		{
			name: "Token from Cookie",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.AddCookie(&http.Cookie{
					Name:  "auth_token",
					Value: "test-cookie-token",
				})
				return req
			},
			wantToken: "test-cookie-token",
			wantErr:   false,
		},
		{
			name: "No token",
			setupReq: func() *http.Request {
				return httptest.NewRequest("GET", "/", nil)
			},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Invalid Authorization header format",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "InvalidFormat")
				return req
			},
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setupReq()

			gotToken, err := getToken(req)

			if (err != nil) != tt.wantErr {
				t.Errorf("getToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotToken != tt.wantToken {
				t.Errorf("getToken() = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}

// Тест для AuthMiddleware
func TestAuthMiddleware(t *testing.T) {
	// Тестовый секретный ключ
	secretKey := "test-secret-key"

	// Создаем валидный токен
	validClaims := &Claims{
		Login: "test-user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims)
	validTokenString, _ := validToken.SignedString([]byte(secretKey))

	// Создаем просроченный токен
	expiredClaims := &Claims{
		Login: "test-user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}

	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, _ := expiredToken.SignedString([]byte(secretKey))

	tests := []struct {
		name           string
		setupReq       func() *http.Request
		wantStatusCode int
		wantLogin      string
	}{
		{
			name: "Valid token in Authorization header",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "Bearer "+validTokenString)
				return req
			},
			wantStatusCode: http.StatusOK,
			wantLogin:      "test-user",
		},
		{
			name: "Valid token in Cookie",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.AddCookie(&http.Cookie{
					Name:  "auth_token",
					Value: validTokenString,
				})
				return req
			},
			wantStatusCode: http.StatusOK,
			wantLogin:      "test-user",
		},
		{
			name: "Expired token",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("Authorization", "Bearer "+expiredTokenString)
				return req
			},
			wantStatusCode: http.StatusUnauthorized,
			wantLogin:      "",
		},
		{
			name: "No token",
			setupReq: func() *http.Request {
				return httptest.NewRequest("GET", "/", nil)
			},
			wantStatusCode: http.StatusUnauthorized,
			wantLogin:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый обработчик, который проверяет наличие логина в контексте
			var capturedLogin string
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				login, ok := r.Context().Value("login").(string)
				if ok {
					capturedLogin = login
				}
				w.WriteHeader(http.StatusOK)
			})

			// Оборачиваем обработчик в middleware для аутентификации
			authHandler := AuthMiddleware(secretKey, testHandler)

			// Создаем тестовый запрос
			req := tt.setupReq()

			// Создаем тестовый ResponseWriter
			recorder := httptest.NewRecorder()

			// Вызываем обработчик
			authHandler.ServeHTTP(recorder, req)

			// Проверяем результаты
			if recorder.Code != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, tt.wantStatusCode)
			}

			if capturedLogin != tt.wantLogin {
				t.Errorf("handler captured wrong login: got %v want %v", capturedLogin, tt.wantLogin)
			}
		})
	}
}
