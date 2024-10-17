package url

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

// IPv4 regex pattern
var ipv4Regex = regexp.MustCompile(`^(([0-9]{1,3}\.){3}[0-9]{1,3})$`)

// Hostname regex pattern for domain names (including TLDs)
var hostnameRegex = regexp.MustCompile(`^(?i:[a-z0-9-]{1,63}\.?)+[a-z]{2,}$`)

// Validates if the given hostname is a valid domain or IP address
func isValidHostname(hostname string) bool {
	if hostname == "" {
		return false
	}

	// Check if the hostname is an IPv4 address
	if ipv4Regex.MatchString(hostname) {
		return true
	}

	// Check if the hostname is a valid domain
	return hostnameRegex.MatchString(hostname)
}

func parseValidURL(rawURL string) (*url.URL, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	// Check protocol (scheme)
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("invalid scheme: %s", parsedURL.Scheme)
	}

	// Check hostname (it should not be empty)
	if parsedURL.Hostname() == "" {
		return nil, fmt.Errorf("invalid hostname")
	}

	// Check hostname (valid domain or IP address)
	if !isValidHostname(parsedURL.Hostname()) {
		return nil, fmt.Errorf("invalid hostname: %s", parsedURL.Hostname())
	}

	return parsedURL, nil
}

func ToHostname(link string) (string, error) {
	parsedURL, err := parseValidURL(link)
	if err != nil {
		return "", err
	}

	return parsedURL.Hostname(), nil
}

// LinkToFullURL maps a relative URL to a full URL based on a provided base URL.
func LinkToFullURL(baseURL, link string) (string, error) {
	// Parse the base URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", errors.New("invalid base URL")
	}

	// Parse the provided link
	u, err := url.Parse(link)
	if err != nil {
		return "", errors.New("invalid link")
	}

	// If the link is already absolute, return it
	if u.IsAbs() {
		return u.String(), nil
	}

	// Otherwise, resolve the relative link against the base URL
	fullURL := base.ResolveReference(u)

	// Return the full URL
	return fullURL.String(), nil
}

func IsHTTPOrHTTPS(link string) bool {
	u, err := url.Parse(link)
	if err != nil {
		return false
	}

	return u.Scheme == "http" || u.Scheme == "https"
}
