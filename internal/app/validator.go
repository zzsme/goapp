package app

import (
	"fmt"
	"regexp"
	"strings"
)

// Validator defines an interface for data validation
type Validator interface {
	Validate(data map[string]interface{}, rules map[string][]string) (bool, map[string]string)
}

// SimpleValidator implements the Validator interface
type SimpleValidator struct{}

// ValidationFunc is a function that validates a value against a rule
type ValidationFunc func(value interface{}, params ...string) bool

// validationRules is a map of rule names to validation functions
var validationRules = map[string]ValidationFunc{
	"required": validateRequired,
	"email":    validateEmail,
	"min":      validateMin,
	"max":      validateMax,
	"in":       validateIn,
	"alpha":    validateAlpha,
	"numeric":  validateNumeric,
}

// CustomValidator is the application's validator
var CustomValidator Validator

// InitValidator initializes the application validator
func InitValidator() {
	CustomValidator = &SimpleValidator{}
	fmt.Println("Validator initialized successfully")
}

// Validate validates data against a set of rules
func (v *SimpleValidator) Validate(data map[string]interface{}, rules map[string][]string) (bool, map[string]string) {
	errors := make(map[string]string)

	for field, fieldRules := range rules {
		value, exists := data[field]

		for _, rule := range fieldRules {
			parts := strings.SplitN(rule, ":", 2)
			ruleName := parts[0]
			ruleParams := []string{}

			if len(parts) > 1 {
				ruleParams = strings.Split(parts[1], ",")
			}

			// Skip non-required fields that don't exist
			if ruleName != "required" && !exists {
				continue
			}

			// Get the validation function
			validationFunc, ok := validationRules[ruleName]
			if !ok {
				errors[field] = fmt.Sprintf("Unknown validation rule: %s", ruleName)
				continue
			}

			// Validate the field
			if !validationFunc(value, ruleParams...) {
				var errorMsg string
				switch ruleName {
				case "required":
					errorMsg = "This field is required."
				case "email":
					errorMsg = "Please enter a valid email address."
				case "min":
					if len(ruleParams) > 0 {
						errorMsg = fmt.Sprintf("This field must be at least %s characters.", ruleParams[0])
					} else {
						errorMsg = "This field does not meet the minimum requirement."
					}
				case "max":
					if len(ruleParams) > 0 {
						errorMsg = fmt.Sprintf("This field must be no more than %s characters.", ruleParams[0])
					} else {
						errorMsg = "This field exceeds the maximum allowed."
					}
				case "in":
					errorMsg = "This field contains an invalid value."
				case "alpha":
					errorMsg = "This field must contain only letters."
				case "numeric":
					errorMsg = "This field must contain only numbers."
				default:
					errorMsg = "This field is invalid."
				}

				errors[field] = errorMsg
				break // Break after the first error for this field
			}
		}
	}

	return len(errors) == 0, errors
}

// validateRequired validates that a field is not empty
func validateRequired(value interface{}, _ ...string) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) != ""
	case int, int64, float64, bool:
		return true
	case []interface{}:
		return len(v) > 0
	case map[string]interface{}:
		return len(v) > 0
	default:
		return false
	}
}

// validateEmail validates that a field is a valid email address
func validateEmail(value interface{}, _ ...string) bool {
	if value == nil {
		return false
	}

	email, ok := value.(string)
	if !ok {
		return false
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// validateMin validates that a field has a minimum length
func validateMin(value interface{}, params ...string) bool {
	if value == nil || len(params) == 0 {
		return false
	}

	min := 0
	fmt.Sscanf(params[0], "%d", &min)

	switch v := value.(type) {
	case string:
		return len(v) >= min
	case int:
		return v >= min
	case int64:
		return int(v) >= min
	case float64:
		return int(v) >= min
	case []interface{}:
		return len(v) >= min
	case map[string]interface{}:
		return len(v) >= min
	default:
		return false
	}
}

// validateMax validates that a field has a maximum length
func validateMax(value interface{}, params ...string) bool {
	if value == nil || len(params) == 0 {
		return false
	}

	max := 0
	fmt.Sscanf(params[0], "%d", &max)

	switch v := value.(type) {
	case string:
		return len(v) <= max
	case int:
		return v <= max
	case int64:
		return int(v) <= max
	case float64:
		return int(v) <= max
	case []interface{}:
		return len(v) <= max
	case map[string]interface{}:
		return len(v) <= max
	default:
		return false
	}
}

// validateIn validates that a field is in a list of values
func validateIn(value interface{}, params ...string) bool {
	if value == nil || len(params) == 0 {
		return false
	}

	strValue := fmt.Sprintf("%v", value)
	for _, param := range params {
		if strValue == param {
			return true
		}
	}

	return false
}

// validateAlpha validates that a field contains only alphabetic characters
func validateAlpha(value interface{}, _ ...string) bool {
	if value == nil {
		return false
	}

	strValue, ok := value.(string)
	if !ok {
		return false
	}

	alphaRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
	return alphaRegex.MatchString(strValue)
}

// validateNumeric validates that a field contains only numeric characters
func validateNumeric(value interface{}, _ ...string) bool {
	if value == nil {
		return false
	}

	switch value.(type) {
	case int, int64, float64:
		return true
	case string:
		strValue := value.(string)
		numericRegex := regexp.MustCompile(`^[0-9]+$`)
		return numericRegex.MatchString(strValue)
	default:
		return false
	}
}
