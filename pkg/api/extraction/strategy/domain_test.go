package strategy

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNewStrategy(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		in := Instruction{
			Selectors: map[string]string{
				"#productTitle":        "text",
				"#priceblock_ourprice": "text",
			},

			OutputFormat: map[string]string{
				"#productTitle":        "string",
				"#priceblock_ourprice": "string",
			},
		}

		userID := uuid.New()

		s, err := NewStrategy(userID, "amazon product prices", in, nil)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if s == nil {
			t.Fatal("unexpected nil strategy")
		}

		if s.ID == uuid.Nil {
			t.Error("unexpected nil strategy ID")
		}

		if s.UserID == uuid.Nil {
			t.Error("unexpected nil user ID")
		}

		if s.Name != "amazon product prices" {
			t.Errorf("unexpected name: %s", s.Name)
		}

		if len(s.Instruction.Selectors) != len(in.Selectors) {
			t.Errorf("unexpected selectors: %v", s.Instruction.Selectors)
		}

		if len(s.Instruction.OutputFormat) != len(in.OutputFormat) {
			t.Errorf("unexpected output format: %v", s.Instruction.OutputFormat)
		}
	})

	t.Run("validates", func(t *testing.T) {
		in := Instruction{}

		userID := uuid.New()

		s, err := NewStrategy(userID, "", in, nil)

		if err == nil {
			t.Fatal("expected error")
		}

		if s != nil {
			t.Fatal("unexpected strategy")
		}

		if !strings.Contains(err.Error(), "name is required") {
			t.Errorf("unexpected error: %v", err)
		}

		if !strings.Contains(err.Error(), "selectors are required") {
			t.Errorf("unexpected error: %v", err)
		}

		if !strings.Contains(err.Error(), "output format is required") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
