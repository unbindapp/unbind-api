package utils

// Makes sure a string contains one interpolation marker, ${} , allow escaping to ignore
func ContainsExactlyOneInterpolationMarker(s string) bool {
	count := 0
	escaped := false

	for i := 0; i < len(s); i++ {
		// Check for escape character
		if s[i] == '\\' {
			escaped = true
			continue
		}

		// Check for ${ sequence when not escaped
		if !escaped && s[i] == '$' && i+2 < len(s) && s[i+1] == '{' && s[i+2] == '}' {
			count++
			i += 2 // Skip the next two characters ('{' and '}')
		}

		escaped = false
	}

	return count == 1
}
