package luna

import (
	"testing"
)

// Тест для LunaAlgorithm
func TestLunaAlgorithm(t *testing.T) {
	tests := []struct {
		name        string
		orderNumber string
		want        bool
	}{
		{
			name:        "Valid credit card number",
			orderNumber: "4532015112830366",
			want:        true,
		},
		{
			name:        "Another valid credit card number",
			orderNumber: "6011000990139424",
			want:        true,
		},
		{
			name:        "Invalid credit card number",
			orderNumber: "4532015112830367",
			want:        false,
		},
		{
			name:        "Non-numeric input",
			orderNumber: "453201511283036a",
			want:        false,
		},
		{
			name:        "Empty input",
			orderNumber: "",
			want:        true, // Сумма 0 делится на 10 без остатка
		},
		{
			name:        "Single digit valid",
			orderNumber: "0",
			want:        true,
		},
		{
			name:        "Single digit invalid",
			orderNumber: "1",
			want:        false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LunaAlgorithm(tt.orderNumber)
			
			if got != tt.want {
				t.Errorf("LunaAlgorithm(%q) = %v, want %v", tt.orderNumber, got, tt.want)
			}
		})
	}
}