package x

import (
	"bytes"
	"encoding/json"
)

func PrettyPrint(b []byte) []byte {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	if err != nil {
		return nil
	}
	return out.Bytes()
}
