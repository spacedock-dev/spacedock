// ABOUTME: Recursive duplicate-member refusal for gate authority JSON documents.
// ABOUTME: Authority is rejected before typed decoding or canonicalization.
package gates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func decodeAuthorityJSON(data []byte, label string, target any) error {
	if err := rejectDuplicateMembers(data); err != nil {
		return fmt.Errorf("%s: %w", label, err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("%s: %w", label, err)
	}
	return nil
}

func rejectDuplicateMembers(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := checkJSONValue(decoder); err != nil {
		return err
	}
	if token, err := decoder.Token(); err != io.EOF {
		if err != nil {
			return err
		}
		return fmt.Errorf("unexpected trailing JSON token %v", token)
	}
	return nil
}

func checkJSONValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		seen := map[string]bool{}
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok {
				return fmt.Errorf("JSON object member name is not a string")
			}
			if seen[key] {
				return fmt.Errorf("duplicate JSON object member %q", key)
			}
			seen[key] = true
			if err := checkJSONValue(decoder); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil {
			return err
		}
		if end != json.Delim('}') {
			return fmt.Errorf("unterminated JSON object")
		}
	case '[':
		for decoder.More() {
			if err := checkJSONValue(decoder); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil {
			return err
		}
		if end != json.Delim(']') {
			return fmt.Errorf("unterminated JSON array")
		}
	default:
		return fmt.Errorf("unexpected JSON delimiter %q", delim)
	}
	return nil
}
