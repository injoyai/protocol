package hostlink

import (
	"encoding/hex"
	"fmt"
	"regexp"
)

func DecodeReadResponse(resp string) (string, error) {
	reg := regexp.MustCompile(`^.{7}(.*).{3}$`)
	matches := reg.FindStringSubmatch(resp)

	if len(matches) <= 1 {
		return "", fmt.Errorf("Wrong response format: %s", resp)
	}

	result, err := hex.DecodeString(matches[1])
	if err != nil {
		return "", fmt.Errorf("Error during decoding response: %s", resp)
	}

	return string(result), nil
}

func DecodeWriteResponse(resp string) (bool, error) {
	result, err := DecodeReadResponse(resp)
	if err != nil {
		return false, err
	}
	return result == "00", nil
}
