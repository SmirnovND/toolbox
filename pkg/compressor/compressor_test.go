package compressor

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Тест для gzipResponseWriter.Write
func TestGzipResponseWriter_Write(t *testing.T) {
	// Создаем тестовый ResponseWriter
	recorder := httptest.NewRecorder()
	
	// Создаем gzip.Writer
	gz := gzip.NewWriter(recorder)
	
	// Создаем gzipResponseWriter
	gzw := gzipResponseWriter{
		ResponseWriter: recorder,
		Writer:         gz,
	}
	
	// Тестовые данные
	testData := []byte("test data")
	
	// Вызываем Write
	n, err := gzw.Write(testData)
	
	// Закрываем gzip.Writer для сброса буфера
	gz.Close()
	
	// Проверяем результаты
	if err != nil {
		t.Errorf("Write() error = %v", err)
	}
	
	if n != len(testData) {
		t.Errorf("Write() n = %v, want %v", n, len(testData))
	}
	
	// Проверяем, что данные были сжаты и записаны
	// Распаковываем данные для проверки
	gzReader, err := gzip.NewReader(bytes.NewReader(recorder.Body.Bytes()))
	if err != nil {
		t.Errorf("Failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()
	
	decompressed, err := io.ReadAll(gzReader)
	if err != nil {
		t.Errorf("Failed to decompress data: %v", err)
	}
	
	if string(decompressed) != string(testData) {
		t.Errorf("Decompressed data = %v, want %v", string(decompressed), string(testData))
	}
}

// Тест для gzipResponseWriter.Header
func TestGzipResponseWriter_Header(t *testing.T) {
	// Создаем тестовый ResponseWriter
	recorder := httptest.NewRecorder()
	
	// Создаем gzip.Writer
	gz := gzip.NewWriter(recorder)
	
	// Создаем gzipResponseWriter
	gzw := gzipResponseWriter{
		ResponseWriter: recorder,
		Writer:         gz,
	}
	
	// Устанавливаем заголовок в оригинальном ResponseWriter
	recorder.Header().Set("Test-Header", "test-value")
	
	// Получаем заголовок через gzipResponseWriter
	header := gzw.Header()
	
	// Проверяем результаты
	if header.Get("Test-Header") != "test-value" {
		t.Errorf("Header().Get(\"Test-Header\") = %v, want %v", header.Get("Test-Header"), "test-value")
	}
}

// Тест для WithCompression
func TestWithCompression(t *testing.T) {
	tests := []struct {
		name            string
		acceptEncoding  string
		wantCompression bool
	}{
		{
			name:            "Client accepts gzip",
			acceptEncoding:  "gzip",
			wantCompression: true,
		},
		{
			name:            "Client accepts gzip with other encodings",
			acceptEncoding:  "deflate, gzip, br",
			wantCompression: true,
		},
		{
			name:            "Client does not accept gzip",
			acceptEncoding:  "deflate, br",
			wantCompression: false,
		},
		{
			name:            "No Accept-Encoding header",
			acceptEncoding:  "",
			wantCompression: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый обработчик
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("test response"))
			})
			
			// Оборачиваем обработчик в middleware для сжатия
			compressedHandler := WithCompression(testHandler)
			
			// Создаем тестовый запрос
			req := httptest.NewRequest("GET", "/", nil)
			if tt.acceptEncoding != "" {
				req.Header.Set("Accept-Encoding", tt.acceptEncoding)
			}
			
			// Создаем тестовый ResponseWriter
			recorder := httptest.NewRecorder()
			
			// Вызываем обработчик
			compressedHandler.ServeHTTP(recorder, req)
			
			// Проверяем результаты
			if tt.wantCompression {
				if recorder.Header().Get("Content-Encoding") != "gzip" {
					t.Errorf("Content-Encoding header = %v, want %v", recorder.Header().Get("Content-Encoding"), "gzip")
				}
				
				// Проверяем, что ответ сжат
				gzReader, err := gzip.NewReader(bytes.NewReader(recorder.Body.Bytes()))
				if err != nil {
					t.Errorf("Failed to create gzip reader: %v", err)
				}
				defer gzReader.Close()
				
				decompressed, err := io.ReadAll(gzReader)
				if err != nil {
					t.Errorf("Failed to decompress response: %v", err)
				}
				
				if string(decompressed) != "test response" {
					t.Errorf("Decompressed response = %v, want %v", string(decompressed), "test response")
				}
			} else {
				if recorder.Header().Get("Content-Encoding") == "gzip" {
					t.Errorf("Content-Encoding header = %v, want empty", recorder.Header().Get("Content-Encoding"))
				}
				
				if recorder.Body.String() != "test response" {
					t.Errorf("Response body = %v, want %v", recorder.Body.String(), "test response")
				}
			}
		})
	}
}

// Тест для WithDecompression
func TestWithDecompression(t *testing.T) {
	tests := []struct {
		name           string
		contentEncoding string
		compressBody   bool
	}{
		{
			name:           "Compressed request",
			contentEncoding: "gzip",
			compressBody:   true,
		},
		{
			name:           "Uncompressed request",
			contentEncoding: "",
			compressBody:   false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый обработчик, который читает тело запроса
			var requestBody string
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := io.ReadAll(r.Body)
				requestBody = string(body)
				w.WriteHeader(http.StatusOK)
			})
			
			// Оборачиваем обработчик в middleware для распаковки
			decompressedHandler := WithDecompression(testHandler)
			
			// Создаем тело запроса
			var reqBody io.Reader
			expectedBody := "test request body"
			
			if tt.compressBody {
				// Сжимаем тело запроса
				var compressedBody bytes.Buffer
				gzWriter := gzip.NewWriter(&compressedBody)
				gzWriter.Write([]byte(expectedBody))
				gzWriter.Close()
				
				reqBody = &compressedBody
			} else {
				reqBody = strings.NewReader(expectedBody)
			}
			
			// Создаем тестовый запрос
			req := httptest.NewRequest("POST", "/", reqBody)
			if tt.contentEncoding != "" {
				req.Header.Set("Content-Encoding", tt.contentEncoding)
			}
			
			// Создаем тестовый ResponseWriter
			recorder := httptest.NewRecorder()
			
			// Вызываем обработчик
			decompressedHandler.ServeHTTP(recorder, req)
			
			// Проверяем результаты
			if requestBody != expectedBody {
				t.Errorf("Request body read by handler = %v, want %v", requestBody, expectedBody)
			}
		})
	}
}