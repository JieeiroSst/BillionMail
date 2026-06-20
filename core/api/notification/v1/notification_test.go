package v1

import (
	"encoding/json"
	"testing"
)

func TestByteArrayMarshalJSON(t *testing.T) {
	b := ByteArray{72, 101, 108, 108, 111}

	data, err := json.Marshal(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := string(data); got != "[72,101,108,108,111]" {
		t.Fatalf("unexpected json: %s", got)
	}
}

func TestByteArrayUnmarshalJSON(t *testing.T) {
	var b ByteArray

	if err := json.Unmarshal([]byte("[72,101,108,108,111]"), &b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(b) != "Hello" {
		t.Fatalf("unexpected result: %v", []byte(b))
	}
}

func TestByteArrayUnmarshalJSONOutOfRange(t *testing.T) {
	var b ByteArray

	if err := json.Unmarshal([]byte("[256]"), &b); err == nil {
		t.Fatalf("expected error for out-of-range value, got nil")
	}
}

func TestSendMobileReqRoundTrip(t *testing.T) {
	req := SendMobileReq{
		Target: "device-token",
		Title:  "Hello",
		Body:   ByteArray("Hi"),
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded SendMobileReq
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(decoded.Body) != "Hi" {
		t.Fatalf("unexpected body after round trip: %v", []byte(decoded.Body))
	}
}
