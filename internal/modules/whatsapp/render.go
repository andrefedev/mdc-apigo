package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// DecodeStrict rejects unknown fields from external payloads.
func decodeJSONStrict(raw []byte, out any) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(out); err != nil {
		return err
	}

	if err := decoder.Decode(new(struct{})); err != io.EOF {
		return fmt.Errorf("expected a single JSON object")
	}

	return nil
}

func decodeRequestError(statusCode int, body []byte) error {
	var env graphErrorEnvelope
	if err := decodeJSONStrict(body, &env); err == nil && env.Error.Message != "" {
		return &RequestError{
			StatusCode: statusCode,
			Code:       env.Error.Code,
			Message:    env.Error.Message,
			Type:       env.Error.Type,
			Subcode:    env.Error.ErrorSubcode,
			TraceID:    env.Error.FBTraceID,
			Body:       string(body),
		}
	}

	return &RequestError{
		Body:       string(body),
		StatusCode: statusCode,
	}
}
