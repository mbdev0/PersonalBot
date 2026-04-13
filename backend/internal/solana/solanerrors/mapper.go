package solanerrors

import "fmt"

func MapTransactionError(raw any) string {
	if raw == nil {
		return ""
	}
	switch err := raw.(type) {
	case string:
		return lookupInstructionError(err)
	case map[string]any:
		return parseErrorMap(err)
	default:
		return fmt.Sprintf("%v", raw)
	}
}

// parseErrorMap handles the map-shaped error from the RPC.
func parseErrorMap(m map[string]any) string {
	if msg := parseInstructionError(m); msg != "" {
		return msg
	}
	// Top-level named errors: e.g. {"BlockhashNotFound": {}}
	for key := range m {
		if msg, ok := solanaInstructionErrors[key]; ok {
			return msg
		}
	}
	return fmt.Sprintf("%v", m)
}

// parseInstructionError handles the InstructionError key, which wraps either
// a custom program error code or a named instruction error.
func parseInstructionError(m map[string]any) string {
	parts, ok := m["InstructionError"].([]any)
	if !ok || len(parts) != 2 {
		return ""
	}

	detail := parts[1]

	if customMap, ok := detail.(map[string]any); ok {
		// Custom program error: {"Custom": 6002}
		if code, ok := customMap["Custom"]; ok {
			return lookupPumpFunError(toInt(code))
		}
	}

	if name, ok := detail.(string); ok {
		// Named instruction error: "InsufficientFundsForFee"
		return lookupInstructionError(name)
	}

	return ""
}

func lookupPumpFunError(code int) string {
	if msg, ok := pumpFunErrors[code]; ok {
		return msg
	}
	return fmt.Sprintf("Program error (code %d)", code)
}

func lookupInstructionError(name string) string {
	if msg, ok := solanaInstructionErrors[name]; ok {
		return msg
	}
	return name
}

func toInt(v any) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return int(val)
	case int64:
		return int(val)
	case uint64:
		return int(val)
	}
	return -1
}
