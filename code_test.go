package spec_test

import (
	"testing"

	spec "github.com/benaskins/axon-spec"
)

func TestBuiltInCodes(t *testing.T) {
	tests := []struct {
		code spec.Code
		want string
	}{
		{spec.MustBePresent, "must-be-present"},
		{spec.MustNotBeEmpty, "must-not-be-empty"},
		{spec.MustBePositive, "must-be-positive"},
	}

	for _, tt := range tests {
		if string(tt.code) != tt.want {
			t.Errorf("got %q, want %q", tt.code, tt.want)
		}
	}
}

func TestCustomCode(t *testing.T) {
	const DebitsMustEqualCredits spec.Code = "debits-must-equal-credits"

	if DebitsMustEqualCredits != "debits-must-equal-credits" {
		t.Fatal("custom code should hold its string value")
	}
}
