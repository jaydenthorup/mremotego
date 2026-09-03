package secrets

import "testing"

func TestParseBitwardenReference(t *testing.T) {
	tests := []struct {
		reference string
		wantID    string
		wantField string
		wantErr   bool
	}{
		{reference: "bw://abc-123", wantID: "abc-123", wantField: "password"},
		{reference: "bw://abc-123/password", wantID: "abc-123", wantField: "password"},
		{reference: "bw://abc-123/username", wantID: "abc-123", wantField: "username"},
		{reference: "bw://abc-123/totp", wantID: "abc-123", wantField: "totp"},
		{reference: "bw://abc-123/notes", wantID: "abc-123", wantField: "notes"},
		{reference: "bw://abc-123/secret", wantErr: true},
		{reference: "bw://abc-123/", wantErr: true},
		{reference: "bw://abc-123/password/extra", wantErr: true},
		{reference: "bw://", wantErr: true},
		{reference: "op://Private/Server/password", wantErr: true},
		{reference: "hunter2", wantErr: true},
		{reference: "", wantErr: true},
	}

	for _, tt := range tests {
		id, field, err := parseBitwardenReference(tt.reference)

		if tt.wantErr {
			if err == nil {
				t.Errorf("parseBitwardenReference(%q) = (%q, %q), want an error", tt.reference, id, field)
			}
			continue
		}

		if err != nil {
			t.Errorf("parseBitwardenReference(%q) returned error: %v", tt.reference, err)
			continue
		}
		if id != tt.wantID || field != tt.wantField {
			t.Errorf("parseBitwardenReference(%q) = (%q, %q), want (%q, %q)",
				tt.reference, id, field, tt.wantID, tt.wantField)
		}
	}
}

func TestBitwardenReference(t *testing.T) {
	if got := bitwardenReference("abc-123"); got != "bw://abc-123" {
		t.Errorf("bitwardenReference = %q, want %q", got, "bw://abc-123")
	}
}

func TestBitwardenReferenceRoundTrip(t *testing.T) {
	id, field, err := parseBitwardenReference(bitwardenReference("abc-123"))
	if err != nil {
		t.Fatalf("round trip returned error: %v", err)
	}
	if id != "abc-123" || field != bitwardenFieldPassword {
		t.Errorf("round trip = (%q, %q), want (abc-123, password)", id, field)
	}
}
