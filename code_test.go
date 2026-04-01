package rule_test

import (
	"testing"

	"github.com/benaskins/axon-rule"
)

func TestBuiltInCodes(t *testing.T) {
	tests := []struct {
		code rule.Code
		want string
	}{
		{rule.MustBePresent, "must-be-present"},
		{rule.MustNotBeEmpty, "must-not-be-empty"},
		{rule.MustBePositive, "must-be-positive"},
	}

	for _, tt := range tests {
		if string(tt.code) != tt.want {
			t.Errorf("got %q, want %q", tt.code, tt.want)
		}
	}
}

func TestCustomCode(t *testing.T) {
	const DebitsMustEqualCredits rule.Code = "debits-must-equal-credits"

	if DebitsMustEqualCredits != "debits-must-equal-credits" {
		t.Fatal("custom code should hold its string value")
	}
}
