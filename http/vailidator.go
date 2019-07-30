package http

// Validator represents a structure that is used to compare string keys
type Validator struct {
	secret string
}

// NewValidator returns a new validator instance with a secret key to
// perform validation against
func NewValidator(secret string) *Validator {
	return &Validator{secret: secret}
}

// Validate compares the value of the secret key stored in the validator
// with the passed in key and determines if they are equivalent
func (v *Validator) Validate(key string) bool {
	return v.secret == key
}
