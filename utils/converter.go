package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
)

var ErrInvalidPreciseAmount = errors.New("invalid precise_amount value")

// ConvertPreciseAmount parses precise_amount strings into arbitrarily large numbers.
// Returns the possibly updated body, whether a conversion happened, and info for logging.
func ConvertPreciseAmount(body []byte) ([]byte, bool, string, *big.Int, error) {
	var jsonBody map[string]interface{}
	if err := json.Unmarshal(body, &jsonBody); err != nil {
		return body, false, "", nil, fmt.Errorf("procesando JSON: %w", err)
	}

	preciseAmount, exists := jsonBody["precise_amount"]
	if !exists {
		return body, false, "", nil, nil
	}

	strValue, ok := preciseAmount.(string)
	if !ok {
		return body, false, "", nil, nil
	}

	intValue, ok := new(big.Int).SetString(strValue, 10)
	if !ok {
		return body, false, strValue, nil, fmt.Errorf("%w: %s", ErrInvalidPreciseAmount, strValue)
	}

	jsonBody["precise_amount"] = json.Number(intValue.String())
	updatedBody, err := json.Marshal(jsonBody)
	if err != nil {
		return body, false, "", nil, fmt.Errorf("convirtiendo JSON: %w", err)
	}

	return updatedBody, true, strValue, intValue, nil
}
