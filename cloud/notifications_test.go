package cloud

import "testing"

func TestMustParse(t *testing.T) {
	result := mustParse("20230618T130047")
	if result.Year() != 2023 {
		t.Errorf("Expected 2023, got %d", result.Year())
	}
	if result.Month() != 6 {
		t.Errorf("Expected 6, got %d", result.Month())
	}
	if result.Day() != 18 {
		t.Errorf("Expected 18, got %d", result.Day())
	}
	if result.Hour() != 13 {
		t.Errorf("Expected 13, got %d", result.Hour())
	}
	if result.Minute() != 0 {
		t.Errorf("Expected 0, got %d", result.Minute())
	}
	if result.Second() != 47 {
		t.Errorf("Expected 47, got %d", result.Second())
	}
}
