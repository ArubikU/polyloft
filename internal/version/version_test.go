package version

import "testing"

func TestString_NotEmpty(t *testing.T) {
	if s := String(); len(s) == 0 {
		t.Fatalf("expected non-empty version string")
	}
}
