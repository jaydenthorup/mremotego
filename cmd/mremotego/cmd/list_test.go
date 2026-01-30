package cmd

import (
"strings"
"testing"

"github.com/jaydenthorup/mremotego/pkg/models"
)

func TestListCommand_Help(t *testing.T) {
if listCmd.Use != "list" {
t.Errorf("Expected Use 'list', got %q", listCmd.Use)
}

if !strings.Contains(listCmd.Short, "List") {
t.Errorf("Expected Short description to contain 'List', got %q", listCmd.Short)
}
}

func TestGetProtocolIcon_ReturnsNonEmpty(t *testing.T) {
protocols := []models.Protocol{
models.ProtocolSSH,
models.ProtocolRDP,
models.ProtocolVNC,
}

for _, proto := range protocols {
if getProtocolIcon(proto) == "" {
t.Errorf("Icon for %s is empty", proto)
}
}
}
