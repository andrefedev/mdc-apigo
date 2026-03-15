package whatsapp

import (
	"bytes"
	"encoding/json"
	"io"
)

// DecodeStrict rejects unknown fields from external payloads.
func decodeStrict(data []byte, out any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	return decoder.Decode(out)
}

// DecodeReaderStrict decodes JSON from a reader rejecting unknown fields.
func decodeReaderStrict(reader io.Reader, out any) error {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	return decoder.Decode(out)
}

// Encode marshals a value as JSON.
func encode(in any) ([]byte, error) {
	return json.Marshal(in)
}
