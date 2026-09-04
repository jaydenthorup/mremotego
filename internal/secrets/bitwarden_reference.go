package secrets

import (
	"fmt"
	"strings"
)

const bitwardenPrefix = SchemeBitwarden + "://"

// Fields of a Bitwarden login item that a reference may point at.
const (
	bitwardenFieldPassword = "password"
	bitwardenFieldUsername = "username"
	bitwardenFieldTotp     = "totp"
	bitwardenFieldNotes    = "notes"
)

var bitwardenFields = map[string]bool{
	bitwardenFieldPassword: true,
	bitwardenFieldUsername: true,
	bitwardenFieldTotp:     true,
	bitwardenFieldNotes:    true,
}

// parseBitwardenReference splits a reference into the item id and the field to
// read. Accepted formats:
//
//	bw://<item-id>            -> the item's password
//	bw://<item-id>/<field>    -> password, username, totp or notes
func parseBitwardenReference(reference string) (id string, field string, err error) {
	if !strings.HasPrefix(reference, bitwardenPrefix) {
		return "", "", fmt.Errorf("reference must start with %s", bitwardenPrefix)
	}

	rest := strings.TrimPrefix(reference, bitwardenPrefix)
	parts := strings.Split(rest, "/")

	switch len(parts) {
	case 1:
		id, field = parts[0], bitwardenFieldPassword
	case 2:
		id, field = parts[0], parts[1]
	default:
		return "", "", fmt.Errorf("reference must be in format %s<item-id>[/<field>]", bitwardenPrefix)
	}

	if id == "" {
		return "", "", fmt.Errorf("reference is missing an item id")
	}

	if !bitwardenFields[field] {
		return "", "", fmt.Errorf("unsupported field %q, expected one of password, username, totp, notes", field)
	}

	return id, field, nil
}

// bitwardenReference builds the reference stored in the configuration for an
// item id.
func bitwardenReference(id string) string {
	return bitwardenPrefix + id
}
