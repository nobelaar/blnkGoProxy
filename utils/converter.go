package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

var ErrInvalidPreciseAmount = errors.New("invalid precise_amount value")

// ConvertPreciseAmount tries to turn a string precise_amount into an int.
// Returns the possibly updated body, whether a conversion happened, and info for logging.
func ConvertPreciseAmount(body []byte) ([]byte, bool, string, int, error) {
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return body, false, "", 0, fmt.Errorf("procesando JSON: %w", err)
	}

	preciseAmount, exists := jsonBody["precise_amount"]
	if !exists {
		return body, false, "", 0, nil
	}

	strValue, ok := preciseAmount.(string)
	if !ok {
		return body, false, "", 0, nil
	}

	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return body, false, strValue, 0, fmt.Errorf("%w: %s", ErrInvalidPreciseAmount, strValue)
	}

	jsonBody["precise_amount"] = intValue
	updatedBody, err := json.Marshal(jsonBody)
	if err != nil {
		return body, false, "", 0, fmt.Errorf("convirtiendo JSON: %w", err)
	}

	return updatedBody, true, strValue, intValue, nil
}
