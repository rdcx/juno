package url

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

func TestLinkToFullURL(t *testing.T) {
	tests := []struct {
		baseURL   string
		link      string
		expected  string
		shouldErr bool
	}{
		// Absolute link, should return the same link
		{
			baseURL:   "https://example.com/subdir/",
			link:      "https://otherdomain.com/resource",
			expected:  "https://otherdomain.com/resource",
			shouldErr: false,
		},
		// Relative link, should be resolved against the base URL
		{
			baseURL:   "https://example.com/subdir/",
			link:      "path/to/resource",
			expected:  "https://example.com/subdir/path/to/resource",
			shouldErr: false,
		},
		// Relative link starting with a slash, should be resolved to root
		{
			baseURL:   "https://example.com/subdir/",
			link:      "/path/to/resource",
			expected:  "https://example.com/path/to/resource",
			shouldErr: false,
		},
		// Relative link with "..", should be resolved to parent directory
		{
			baseURL:   "https://example.com/subdir/",
			link:      "../resource",
			expected:  "https://example.com/resource",
			shouldErr: false,
		},
		// Invalid base URL, should return an error
		{
			baseURL:   "://invalid-url",
			link:      "path/to/resource",
			expected:  "",
			shouldErr: true,
		},
		// Invalid relative link, should return an error
		{
			baseURL:   "https://example.com/subdir/",
			link:      "://invalid-link",
			expected:  "",
			shouldErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.link, func(t *testing.T) {
			result, err := LinkToFullURL(test.baseURL, test.link)

			if test.shouldErr && err == nil {
				t.Fatalf("expected error but got none")
			}

			if !test.shouldErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != test.expected {
				t.Errorf("expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestIsHTTPOrHTTPS(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "http",
			url:      "http://example.com",
			expected: true,
		},
		{
			name:     "https",
			url:      "https://example.com",
			expected: true,
		},
		{
			name:     "ftp",
			url:      "ftp://example.com",
			expected: false,
		},
		{
			name:     "mailto",
			url:      "mailto:",
			expected: false,
		},
		{
			name:     "tel",
			url:      "tel:",
			expected: false,
		},
		{
			name:     "data",
			url:      "data:",
			expected: false,
		},
		{
			name:     "javascript",
			url:      "javascript:",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsHTTPOrHTTPS(tt.url)

			if got != tt.expected {
				t.Errorf("IsHTTPOrHTTPS() = %v, want %v", got, tt.expected)
			}
		})
	}
}
