package paramsparser

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// Тестовая структура для JSONParse
type TestStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Тест для JSONParse
func TestJSONParse(t *testing.T) {
	tests := []struct {
		name    string
		body    interface{}
		want    *TestStruct
		wantErr bool
	}{
		{
			name: "Valid JSON",
			body: TestStruct{
				Name: "John",
				Age:  30,
			},
			want: &TestStruct{
				Name: "John",
				Age:  30,
			},
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			body:    "invalid json",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty JSON",
			body:    TestStruct{},
			want:    &TestStruct{},
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Подготовка тела запроса
			var reqBody []byte
			var err error
			
			switch v := tt.body.(type) {
			case string:
				reqBody = []byte(v)
			default:
				reqBody, err = json.Marshal(v)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}
			
			// Создаем тестовый запрос
			req := httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			
			// Создаем тестовый ResponseWriter
			recorder := httptest.NewRecorder()
			
			// Вызываем JSONParse
			got, err := JSONParse[TestStruct](recorder, req)
			
			// Проверяем результаты
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONParse() = %v, want %v", got, tt.want)
			}
			
			// Проверяем статус ответа в случае ошибки
			if tt.wantErr && recorder.Code != http.StatusBadRequest {
				t.Errorf("JSONParse() status code = %v, want %v", recorder.Code, http.StatusBadRequest)
			}
		})
	}
}

// Тест для TextParse
func TestTextParse(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    string
		wantErr bool
	}{
		{
			name:    "Valid text",
			body:    "test text",
			want:    "test text",
			wantErr: false,
		},
		{
			name:    "Empty text",
			body:    "",
			want:    "",
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый запрос
			req := httptest.NewRequest("POST", "/", strings.NewReader(tt.body))
			
			// Создаем тестовый ResponseWriter
			recorder := httptest.NewRecorder()
			
			// Вызываем TextParse
			got, err := TextParse(recorder, req)
			
			// Проверяем результаты
			if (err != nil) != tt.wantErr {
				t.Errorf("TextParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if got != tt.want {
				t.Errorf("TextParse() = %v, want %v", got, tt.want)
			}
		})
	}
}