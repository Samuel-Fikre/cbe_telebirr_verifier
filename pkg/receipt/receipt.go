package receipt

// Receipt represents a telebirr transaction receipt
type Receipt struct {
	ParsedFields     map[string]interface{}
	PreDefinedFields map[string]interface{}
}

// New creates a new Receipt instance
func New(parsedFields, preDefinedFields map[string]interface{}) *Receipt {
	return &Receipt{
		ParsedFields:     parsedFields,
		PreDefinedFields: preDefinedFields,
	}
}

// Equals compares two values for equality
func (r *Receipt) Equals(a, b interface{}) bool {
	switch v1 := a.(type) {
	case string:
		if v2, ok := b.(string); ok {
			return v1 == v2
		}
	case float64:
		switch v2 := b.(type) {
		case float64:
			return v1 == v2
		case int:
			return v1 == float64(v2)
		}
	case int:
		switch v2 := b.(type) {
		case float64:
			return float64(v1) == v2
		case int:
			return v1 == v2
		}
	}
	return false
}

// Verify checks parsed fields against predefined fields using a custom verification function
func (r *Receipt) Verify(callback func(map[string]interface{}, map[string]interface{}) bool) bool {
	return callback(r.ParsedFields, r.PreDefinedFields)
}

// VerifyAll checks all parsed fields against their corresponding predefined fields
func (r *Receipt) VerifyAll(doNotCompare []string) bool {
	if len(r.ParsedFields) == 0 {
		return false
	}

	for key, preDefinedValue := range r.PreDefinedFields {
		if contains(doNotCompare, key) {
			continue
		}
		if parsedValue, exists := r.ParsedFields[key]; !exists || !r.Equals(parsedValue, preDefinedValue) {
			return false
		}
	}
	return true
}

// VerifyOnly checks only the specified fields against their corresponding predefined fields
func (r *Receipt) VerifyOnly(fieldNames []string) bool {
	if len(fieldNames) == 0 {
		return false
	}

	for _, key := range fieldNames {
		parsedValue, parsedExists := r.ParsedFields[key]
		preDefinedValue, preDefinedExists := r.PreDefinedFields[key]

		if !parsedExists || !preDefinedExists || !r.Equals(parsedValue, preDefinedValue) {
			return false
		}
	}
	return true
}

// Helper function to check if a string is in a slice
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
