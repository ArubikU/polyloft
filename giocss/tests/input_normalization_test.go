package giocss_test

import (
	"testing"

	giocss "github.com/ArubikU/giocss"
)

func TestInputExternalValueSupportsNumericValue(t *testing.T) {
	props := map[string]any{"value": -1}
	got := giocss.InputExternalValue(props)
	if got != "-1" {
		t.Fatalf("InputExternalValue() = %q, want %q", got, "-1")
	}
}

func TestNormalizeNumberInputFinalizeClampsAndRespectsStep(t *testing.T) {
	props := map[string]any{"min": 0, "max": 10, "step": 0.5}
	got := giocss.NormalizeNumberInput("10.7", props, true)
	if got != "10" {
		t.Fatalf("NormalizeNumberInput() = %q, want %q", got, "10")
	}
}

func TestNormalizeNumberInputLiveAllowsPartial(t *testing.T) {
	got := giocss.NormalizeNumberInput("-", map[string]any{}, false)
	if got != "-" {
		t.Fatalf("NormalizeNumberInput(live) = %q, want %q", got, "-")
	}
}

func TestNormalizeDateInputFinalize(t *testing.T) {
	got := giocss.NormalizeDateInput("2026/03/17", true)
	if got != "2026-03-17" {
		t.Fatalf("NormalizeDateInput() = %q, want %q", got, "2026-03-17")
	}
}

func TestNormalizeDateInputRejectsInvalid(t *testing.T) {
	got := giocss.NormalizeDateInput("2026-99-99", true)
	if got != "" {
		t.Fatalf("NormalizeDateInput(invalid) = %q, want empty", got)
	}
}

func TestNormalizeTimeInputFinalizeWithSecondsStep(t *testing.T) {
	props := map[string]any{"step": 1}
	got := giocss.NormalizeTimeInput("02:30", props, true)
	if got != "02:30:00" {
		t.Fatalf("NormalizeTimeInput() = %q, want %q", got, "02:30:00")
	}
}

func TestNormalizeTimeInputRejectsInvalid(t *testing.T) {
	got := giocss.NormalizeTimeInput("25:99", map[string]any{}, true)
	if got != "" {
		t.Fatalf("NormalizeTimeInput(invalid) = %q, want empty", got)
	}
}
