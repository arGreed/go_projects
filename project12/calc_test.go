package calculator

import (
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		first     float64
		second    float64
		operand   string
		expected  float64
		expectErr bool
	}{
		{10, 5, "+", 15, false},
		{10, 5, "-", 5, false},
		{10, 5, "*", 50, false},
		{10, 5, "/", 2, false},
		{10, 0, "/", 0, true}, // деление на ноль
		{10, 5, "%", 0, true}, // несуществующая операция
	}

	for _, test := range tests {
		result, err := calculate(&test.first, &test.second, &test.operand)
		if (err != nil) != test.expectErr {
			t.Errorf("calculate(%v, %v, %v) = error: %v, want error: %v", test.first, test.second, test.operand, err, test.expectErr)
		}
		if result != test.expected {
			t.Errorf("calculate(%v, %v, %v) = %v, want %v", test.first, test.second, test.operand, result, test.expected)
		}
	}
}
