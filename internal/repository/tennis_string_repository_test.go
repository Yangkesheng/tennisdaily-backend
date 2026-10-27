package repository

import "testing"

func TestStringSearchPattern(t *testing.T) {
	if got := stringSearchPattern("50%_!"); got != "%50!%!_!!%" {
		t.Fatalf("LIKE wildcard escaping failed: %q", got)
	}
}
