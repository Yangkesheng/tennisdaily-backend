package service

import "testing"

func TestDeriveBrandSlug(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{name: "Asics", want: "asics"},
		{name: "adidas", want: "adidas"},
		{name: "New Balance", want: "new-balance"},
		{name: "Wilson Pro Staff", want: "wilson-pro-staff"},
		{name: "  Yonex  ", want: "yonex"},
		{name: "李宁", want: ""},
	}
	for _, c := range cases {
		slug := deriveBrandSlug(c.name)
		if c.want != "" {
			if slug != c.want {
				t.Fatalf("deriveBrandSlug(%q) = %q, want %q", c.name, slug, c.want)
			}
			continue
		}
		if len(slug) < 7 || slug[:6] != "brand-" {
			t.Fatalf("deriveBrandSlug(%q) = %q, want brand- prefix", c.name, slug)
		}
	}
}

func TestIsAdminUser(t *testing.T) {
	ids := []int64{1, 7}
	if !IsAdminUser(1, ids) {
		t.Fatal("expected user 1 to be admin")
	}
	if !IsAdminUser(7, ids) {
		t.Fatal("expected user 7 to be admin")
	}
	if IsAdminUser(2, ids) {
		t.Fatal("expected user 2 not to be admin")
	}
	if IsAdminUser(0, ids) {
		t.Fatal("expected user 0 not to be admin")
	}
	if IsAdminUser(1, nil) {
		t.Fatal("expected no admins when whitelist empty")
	}
}
