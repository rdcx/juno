package link

import "testing"

func TestToHostname(t *testing.T) {

	tests := []struct {
		name     string
		link     string
		expected string
	}{
		{
			name:     "simple",
			link:     "http://example.com",
			expected: "example.com",
		},
		{
			name:     "with path",
			link:     "http://example.net/path",
			expected: "example.net",
		},
		{
			name:     "with query",
			link:     "http://shop.org/path?query=1",
			expected: "shop.org",
		},
		{
			name:     "with subdomain",
			link:     "http://sub.clothes-shop.org/path?query=1",
			expected: "sub.clothes-shop.org",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToHostname(tt.link)

			if err != nil {
				t.Errorf("ToHostname() error = %v, want nil", err)
			}

			if got != tt.expected {
				t.Errorf("ToHostname() = %v, want %v", got, tt.expected)
			}
		})
	}
}
