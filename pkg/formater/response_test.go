package formater

import (
	"encoding/json"
	"reflect"
	"testing"
)

// Тест для JSONResponse
func TestJSONResponse(t *testing.T) {
	tests := []struct {
		name    string
		data    interface{}
		want    []byte
		wantErr bool
	}{
		{
			name: "Simple struct",
			data: struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}{
				Name: "John",
				Age:  30,
			},
			want:    []byte(`{"name":"John","age":30}`),
			wantErr: false,
		},
		{
			name:    "Simple map",
			data:    map[string]string{"key": "value"},
			want:    []byte(`{"key":"value"}`),
			wantErr: false,
		},
		{
			name:    "Simple slice",
			data:    []int{1, 2, 3},
			want:    []byte(`[1,2,3]`),
			wantErr: false,
		},
		{
			name:    "Simple string",
			data:    "test",
			want:    []byte(`"test"`),
			wantErr: false,
		},
		{
			name:    "Simple int",
			data:    42,
			want:    []byte(`42`),
			wantErr: false,
		},
		{
			name:    "Nil",
			data:    nil,
			want:    []byte(`null`),
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONResponse(tt.data)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			// Сравниваем JSON-структуры, а не строки, чтобы избежать проблем с форматированием
			var gotJSON, wantJSON interface{}
			
			if err := json.Unmarshal(got, &gotJSON); err != nil {
				t.Errorf("Failed to unmarshal result JSON: %v", err)
				return
			}
			
			if err := json.Unmarshal(tt.want, &wantJSON); err != nil {
				t.Errorf("Failed to unmarshal expected JSON: %v", err)
				return
			}
			
			if !reflect.DeepEqual(gotJSON, wantJSON) {
				t.Errorf("JSONResponse() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}