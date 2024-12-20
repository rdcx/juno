package util

import (
	"fmt"
	"testing"
)

func TestBcryptPassword(t *testing.T) {
	password := "password"
	hash, err := BcryptPassword(password)
	if err != nil {
		t.Errorf("BcryptPassword() error = %v, want nil", err)
	}

	if err := CompareBcryptPassword(hash, password); err != nil {
		t.Errorf("CompareBcryptPassword() error = %v, want nil", err)
	}
}

func TestValidationErrs(t *testing.T) {
	errs := []error{
		fmt.Errorf("error1"),
		fmt.Errorf("error2"),
	}

	ve := ValidationErrs(errs)

	if ve.Error() != "validation errors: error1, error2" {
		t.Errorf("ValidationErrs() = %s, want validation errors: error1, error2", ve.Error())
	}
}

func TestIsValidHostname(t *testing.T) {

	tests := []struct {
		name     string
		hostname string
		want     bool
	}{
		{
			name:     "valid hostname",
			hostname: "example.com",
			want:     true,
		},
		{
			name:     "valid hostname with subdomain",
			hostname: "sub.example.com",
			want:     true,
		},
		{
			name:     "invalid hostname",
			hostname: "exam!ple",
			want:     false,
		},
		{
			name:     "invalid hostname with subdomain",
			hostname: "sub.exam!ple.com",
			want:     false,
		},
		{
			name:     "valid IP address",
			hostname: "10.0.0.1",
			want:     true,
		},
		{
			name:     "invalid IP address",
			hostname: "10.0.0.256",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidHostname(tt.hostname)
			if got != tt.want {
				t.Errorf("IsValidHostname() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidPort(t *testing.T) {

	if IsValidPort("0") {
		t.Errorf("IsValidPort(0) = true, want false")
	}

	for i := 1; i < 65536; i++ {
		if !IsValidPort(fmt.Sprintf("%d", i)) {
			t.Errorf("IsValidPort(%d) = false, want true", i)
		}
	}

	if IsValidPort("-1") {
		t.Errorf("IsValidPort(-1) = true, want false")
	}

	if IsValidPort("65536") {
		t.Errorf("IsValidPort(65536) = true, want false")
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{
			name:  "valid email",
			email: "valid@email.com",
			want:  true,
		},
		{
			name:  "valid email with subdomain and plus sign",
			email: "valid+email@one.example.com",
			want:  true,
		},
		{
			name:  "invalid email",
			email: "invalid",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidEmail(tt.email)
			if got != tt.want {
				t.Errorf("IsValidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
