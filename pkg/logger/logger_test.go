package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Тест для loggingResponseWriter.Write
func TestLoggingResponseWriter_Write(t *testing.T) {
	// Создаем тестовый ResponseWriter
	recorder := httptest.NewRecorder()
	
	// Создаем responseData
	respData := &responseData{
		status: 0,
		size:   0,
	}
	
	// Создаем loggingResponseWriter
	lw := loggingResponseWriter{
		ResponseWriter: recorder,
		responseData:   respData,
	}
	
	// Тестовые данные
	testData := []byte("test data")
	
	// Вызываем Write
	n, err := lw.Write(testData)
	
	// Проверяем результаты
	if err != nil {
		t.Errorf("Write() error = %v", err)
	}
	
	if n != len(testData) {
		t.Errorf("Write() n = %v, want %v", n, len(testData))
	}
	
	if respData.size != len(testData) {
		t.Errorf("responseData.size = %v, want %v", respData.size, len(testData))
	}
	
	// Проверяем, что данные были записаны в оригинальный ResponseWriter
	if recorder.Body.String() != string(testData) {
		t.Errorf("ResponseWriter.Body = %v, want %v", recorder.Body.String(), string(testData))
	}
}

// Тест для loggingResponseWriter.WriteHeader
func TestLoggingResponseWriter_WriteHeader(t *testing.T) {
	// Создаем тестовый ResponseWriter
	recorder := httptest.NewRecorder()
	
	// Создаем responseData
	respData := &responseData{
		status: 0,
		size:   0,
	}
	
	// Создаем loggingResponseWriter
	lw := loggingResponseWriter{
		ResponseWriter: recorder,
		responseData:   respData,
	}
	
	// Тестовый статус
	testStatus := http.StatusOK
	
	// Вызываем WriteHeader
	lw.WriteHeader(testStatus)
	
	// Проверяем результаты
	if respData.status != testStatus {
		t.Errorf("responseData.status = %v, want %v", respData.status, testStatus)
	}
	
	// Проверяем, что статус был установлен в оригинальном ResponseWriter
	if recorder.Code != testStatus {
		t.Errorf("ResponseWriter.Code = %v, want %v", recorder.Code, testStatus)
	}
}

// Тест для WithLogging
func TestWithLogging(t *testing.T) {
	// Создаем тестовый обработчик
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})
	
	// Оборачиваем обработчик в middleware для логирования
	loggedHandler := WithLogging(testHandler)
	
	// Создаем тестовый запрос
	req := httptest.NewRequest("GET", "/test", nil)
	
	// Создаем тестовый ResponseWriter
	recorder := httptest.NewRecorder()
	
	// Вызываем обработчик
	loggedHandler.ServeHTTP(recorder, req)
	
	// Проверяем результаты
	if recorder.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", recorder.Code, http.StatusOK)
	}
	
	expectedBody := "test response"
	if recorder.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), expectedBody)
	}
}