package giocss_test

import (
	"image/color"
	"testing"

	giocss "github.com/ArubikU/giocss"
)

func TestParseHexColorSupportsShortAndAlphaForms(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  color.NRGBA
	}{
		{name: "short hex", input: "#abc", want: color.NRGBA{R: 0xAA, G: 0xBB, B: 0xCC, A: 0xFF}},
		{name: "hex with alpha", input: "#11223344", want: color.NRGBA{R: 0x11, G: 0x22, B: 0x33, A: 0x44}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := asNRGBA(giocss.ParseHexColor(tt.input, color.NRGBA{}))
			if got != tt.want {
				t.Fatalf("ParseHexColor(%q) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseHexColorExtractsFirstGradientStop(t *testing.T) {
	got := asNRGBA(giocss.ParseHexColor("linear-gradient(90deg, #ff0000 0%, #0000ff 100%)", color.NRGBA{}))
	want := color.NRGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}
	if got != want {
		t.Fatalf("ParseHexColor(gradient) = %#v, want %#v", got, want)
	}
}

func TestParseRGBColorSupportsPercentAndAlpha(t *testing.T) {
	got := asNRGBA(giocss.ParseRGBColor("rgba(100%, 50%, 0%, 0.5)", color.NRGBA{}))
	want := color.NRGBA{R: 0xFF, G: 0x7F, B: 0x00, A: 0x7F}
	if got != want {
		t.Fatalf("ParseRGBColor() = %#v, want %#v", got, want)
	}
}

func TestParseHSLColorSupportsCanonicalGreen(t *testing.T) {
	got := asNRGBA(giocss.ParseHSLColor("hsl(120, 100%, 50%)", color.NRGBA{}))
	want := color.NRGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}
	if got != want {
		t.Fatalf("ParseHSLColor() = %#v, want %#v", got, want)
	}
}

func TestCSSTextTransformCapitalize(t *testing.T) {
	got := giocss.CSSTextTransform("hola mundo polyloft", map[string]string{"text-transform": "capitalize"})
	if got != "Hola Mundo Polyloft" {
		t.Fatalf("CSSTextTransform() = %q, want %q", got, "Hola Mundo Polyloft")
	}
}

func asNRGBA(c color.Color) color.NRGBA {
	return color.NRGBAModel.Convert(c).(color.NRGBA)
}
