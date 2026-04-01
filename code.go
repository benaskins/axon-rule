package rule

// Code is a typed identifier for a violation. Domain packages define their
// own codes as constants. axon-spec provides a small set of common codes.
type Code string

// Common codes for universal business rules.
const (
	MustBePresent  Code = "must-be-present"
	MustNotBeEmpty Code = "must-not-be-empty"
	MustBePositive Code = "must-be-positive"
)
