//go:build !windows

package launcher

// The Windows Credential Manager has no counterpart on other platforms; RDP
// clients there take credentials by other means.

func writeGenericCredential(target, username, password string) error { return nil }

func deleteGenericCredential(target string) error { return nil }
