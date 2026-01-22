package launcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/yourusername/mremotego/internal/secrets"
	"github.com/yourusername/mremotego/pkg/models"
)

// Launcher handles launching connections
type Launcher struct {
	onePasswordProvider *secrets.OnePasswordProvider
}

// NewLauncher creates a new launcher
func NewLauncher() *Launcher {
	return &Launcher{
		onePasswordProvider: secrets.NewOnePasswordProvider(),
	}
}

// hideConsoleWindow sets the command attributes to hide console windows on Windows
func hideConsoleWindow(cmd *exec.Cmd) {
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x08000000, // CREATE_NO_WINDOW
		}
	}
}

// Launch launches a connection based on its protocol
func (l *Launcher) Launch(conn *models.Connection) error {
	if conn.IsFolder() {
		return fmt.Errorf("cannot launch a folder")
	}

	// Resolve 1Password reference if needed (make a copy to avoid modifying the original)
	resolvedConn := *conn
	if l.onePasswordProvider.IsReference(conn.Password) {
		resolvedConn.Password = l.onePasswordProvider.ResolveIfReference(conn.Password)
	}

	switch resolvedConn.Protocol {
	case models.ProtocolSSH:
		return l.launchSSH(&resolvedConn)
	case models.ProtocolRDP:
		return l.launchRDP(&resolvedConn)
	case models.ProtocolVNC:
		return l.launchVNC(&resolvedConn)
	case models.ProtocolHTTP, models.ProtocolHTTPS:
		return l.launchHTTP(&resolvedConn)
	case models.ProtocolTelnet:
		return l.launchTelnet(&resolvedConn)
	default:
		return fmt.Errorf("unsupported protocol: %s", resolvedConn.Protocol)
	}
}

// launchSSH launches an SSH connection using PuTTY on Windows, ssh otherwise
func (l *Launcher) launchSSH(conn *models.Connection) error {
	port := conn.Port
	if port == 0 {
		port = models.ProtocolSSH.GetDefaultPort()
	}

	if runtime.GOOS == "windows" {
		// Use PuTTY on Windows
		args := []string{}

		// Add SSH protocol
		args = append(args, "-ssh")

		// Add port
		args = append(args, "-P", strconv.Itoa(port))

		// Add username if provided
		if conn.Username != "" {
			args = append(args, "-l", conn.Username)
		}

		// Add password if provided (for auto-login)
		if conn.Password != "" {
			args = append(args, "-pw", conn.Password)
		}

		// Add extra args if provided
		if conn.ExtraArgs != "" {
			args = append(args, conn.ExtraArgs)
		}

		// Add hostname last
		args = append(args, conn.Host)

		// Try putty.exe first, fall back to ssh if not found
		cmd := exec.Command("putty.exe", args...)
		hideConsoleWindow(cmd)
		if err := cmd.Start(); err != nil {
			// Fall back to ssh command
			return l.launchSSHFallback(conn)
		}
		return nil
	}

	// Use ssh on Linux/Mac
	return l.launchSSHFallback(conn)
}

// launchSSHFallback uses the standard ssh command as fallback
func (l *Launcher) launchSSHFallback(conn *models.Connection) error {
	port := conn.Port
	if port == 0 {
		port = models.ProtocolSSH.GetDefaultPort()
	}

	args := []string{}

	// Add port
	args = append(args, "-p", strconv.Itoa(port))

	// Add username if provided
	target := conn.Host
	if conn.Username != "" {
		target = conn.Username + "@" + conn.Host
	}

	args = append(args, target)

	// Add extra args if provided
	if conn.ExtraArgs != "" {
		args = append(args, conn.ExtraArgs)
	}

	cmd := exec.Command("ssh", args...)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	hideConsoleWindow(cmd)

	return cmd.Start()
}

// launchRDP launches an RDP connection
func (l *Launcher) launchRDP(conn *models.Connection) error {
	port := conn.Port
	if port == 0 {
		port = models.ProtocolRDP.GetDefaultPort()
	}

	target := conn.Host
	if port != 3389 {
		target = fmt.Sprintf("%s:%d", conn.Host, port)
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Store credentials in Windows Credential Manager if password is provided
		if conn.Username != "" && conn.Password != "" {
			if err := l.storeWindowsCredential(conn); err != nil {
				// Continue anyway - mstsc will prompt if credentials aren't stored
				fmt.Printf("Warning: Failed to store credentials: %v\n", err)
			}
		}

		// Create a temporary .rdp file with connection settings
		rdpFile, err := l.createRDPFile(conn, target)
		if err != nil {
			return fmt.Errorf("failed to create RDP file: %w", err)
		}

		// Launch mstsc with the RDP file
		cmd = exec.Command("mstsc", rdpFile)
		hideConsoleWindow(cmd)

	case "linux", "darwin":
		// Use xfreerdp on Linux/Mac
		args := []string{
			"/v:" + target,
			"/cert:ignore",
		}

		if conn.Username != "" {
			args = append(args, "/u:"+conn.Username)
		}

		if conn.Domain != "" {
			args = append(args, "/d:"+conn.Domain)
		}

		if conn.Password != "" {
			args = append(args, "/p:"+conn.Password)
		}

		if conn.Resolution != "" {
			args = append(args, "/size:"+conn.Resolution)
		} else {
			args = append(args, "/f") // Fullscreen by default
		}

		if conn.ColorDepth > 0 {
			args = append(args, fmt.Sprintf("/bpp:%d", conn.ColorDepth))
		}

		if conn.ExtraArgs != "" {
			args = append(args, conn.ExtraArgs)
		}

		cmd = exec.Command("xfreerdp", args...)
		hideConsoleWindow(cmd)

	default:
		return fmt.Errorf("RDP not supported on this platform")
	}

	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	return cmd.Start()
}

// launchVNC launches a VNC connection
func (l *Launcher) launchVNC(conn *models.Connection) error {
	port := conn.Port
	if port == 0 {
		port = models.ProtocolVNC.GetDefaultPort()
	}

	target := fmt.Sprintf("%s:%d", conn.Host, port)

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Try common VNC clients on Windows
		cmd = exec.Command("vncviewer", target)
	case "linux":
		// Try vncviewer on Linux
		cmd = exec.Command("vncviewer", target)
	case "darwin":
		// Use open with vnc:// protocol on Mac
		target = fmt.Sprintf("vnc://%s:%d", conn.Host, port)
		cmd = exec.Command("open", target)
	default:
		return fmt.Errorf("VNC not supported on this platform")
	}

	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	hideConsoleWindow(cmd)

	return cmd.Start()
}

// launchHTTP launches an HTTP/HTTPS connection in default browser
func (l *Launcher) launchHTTP(conn *models.Connection) error {
	url := fmt.Sprintf("%s://%s", conn.Protocol, conn.Host)

	if conn.Port != 0 {
		url = fmt.Sprintf("%s:%d", url, conn.Port)
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
		hideConsoleWindow(cmd)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("cannot open browser on this platform")
	}

	return cmd.Start()
}

// launchTelnet launches a telnet connection
func (l *Launcher) launchTelnet(conn *models.Connection) error {
	port := conn.Port
	if port == 0 {
		port = models.ProtocolTelnet.GetDefaultPort()
	}

	args := []string{conn.Host, strconv.Itoa(port)}

	cmd := exec.Command("telnet", args...)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	hideConsoleWindow(cmd)

	return cmd.Start()
}

// createRDPFile creates a temporary .rdp file with connection settings
func (l *Launcher) createRDPFile(conn *models.Connection, target string) (string, error) {
	// Create temp directory for RDP files
	tempDir := filepath.Join(os.TempDir(), "mremotego")
	if err := os.MkdirAll(tempDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Sanitize connection name for filename (remove invalid characters)
	safeName := sanitizeFilename(conn.Name)

	// Create RDP file with .rdp extension
	rdpPath := filepath.Join(tempDir, safeName+".rdp")

	// Build RDP file content (UTF-8 encoding)
	rdpContent := fmt.Sprintf("full address:s:%s\r\n", target)

	if conn.Username != "" {
		// Include domain if specified
		if conn.Domain != "" {
			rdpContent += fmt.Sprintf("domain:s:%s\r\n", conn.Domain)
			rdpContent += fmt.Sprintf("username:s:%s\\%s\r\n", conn.Domain, conn.Username)
		} else {
			rdpContent += fmt.Sprintf("username:s:%s\r\n", conn.Username)
		}
	}

	// Note: Password is stored in Windows Credential Manager, not in RDP file
	// This is the same approach that mstsc uses for saved credentials

	// Common RDP settings (must use \r\n line endings for Windows)
	rdpContent += "screen mode id:i:2\r\n" // Fullscreen
	rdpContent += "use multimon:i:0\r\n"
	rdpContent += "session bpp:i:32\r\n"
	rdpContent += "compression:i:1\r\n"
	rdpContent += "keyboardhook:i:2\r\n"
	rdpContent += "audiocapturemode:i:0\r\n"
	rdpContent += "videoplaybackmode:i:1\r\n"
	rdpContent += "connection type:i:7\r\n"
	rdpContent += "networkautodetect:i:1\r\n"
	rdpContent += "bandwidthautodetect:i:1\r\n"
	rdpContent += "displayconnectionbar:i:1\r\n"
	rdpContent += "enableworkspacereconnect:i:0\r\n"
	rdpContent += "disable wallpaper:i:0\r\n"
	rdpContent += "allow font smoothing:i:1\r\n"
	rdpContent += "allow desktop composition:i:1\r\n"
	rdpContent += "disable full window drag:i:0\r\n"
	rdpContent += "disable menu anims:i:0\r\n"
	rdpContent += "disable themes:i:0\r\n"
	rdpContent += "disable cursor setting:i:0\r\n"
	rdpContent += "bitmapcachepersistenable:i:1\r\n"
	rdpContent += "audiomode:i:0\r\n"
	rdpContent += "redirectprinters:i:0\r\n"
	rdpContent += "redirectcomports:i:0\r\n"
	rdpContent += "redirectsmartcards:i:0\r\n"
	rdpContent += "redirectclipboard:i:1\r\n"
	rdpContent += "redirectposdevices:i:0\r\n"
	rdpContent += "autoreconnection enabled:i:1\r\n"
	rdpContent += "authentication level:i:2\r\n"
	rdpContent += "prompt for credentials:i:0\r\n"
	rdpContent += "negotiate security layer:i:1\r\n"
	rdpContent += "remoteapplicationmode:i:0\r\n"
	rdpContent += "alternate shell:s:\r\n"
	rdpContent += "shell working directory:s:\r\n"
	rdpContent += "gatewayhostname:s:\r\n"
	rdpContent += "gatewayusagemethod:i:4\r\n"
	rdpContent += "gatewaycredentialssource:i:4\r\n"
	rdpContent += "gatewayprofileusagemethod:i:0\r\n"
	rdpContent += "promptcredentialonce:i:0\r\n"
	rdpContent += "gatewaybrokeringtype:i:0\r\n"
	rdpContent += "use redirection server name:i:0\r\n"
	rdpContent += "rdgiskdcproxy:i:0\r\n"
	rdpContent += "kdcproxyname:s:\r\n"

	// Resolution settings
	if conn.Resolution != "" {
		// Parse resolution like "1920x1080"
		rdpContent += fmt.Sprintf("desktopwidth:i:1920\r\n")
		rdpContent += fmt.Sprintf("desktopheight:i:1080\r\n")
	} else {
		rdpContent += "smart sizing:i:1\r\n"
	}

	// Color depth
	if conn.ColorDepth > 0 {
		rdpContent += fmt.Sprintf("session bpp:i:%d\r\n", conn.ColorDepth)
	}

	// Write RDP file
	if err := os.WriteFile(rdpPath, []byte(rdpContent), 0600); err != nil {
		return "", fmt.Errorf("failed to write RDP file: %w", err)
	}

	return rdpPath, nil
}

// sanitizeFilename removes or replaces invalid characters from a filename
func sanitizeFilename(name string) string {
	// Replace invalid Windows filename characters with underscore
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := name
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}
	// Trim spaces and dots from start/end
	result = strings.Trim(result, " .")
	// If empty after sanitization, use a default name
	if result == "" {
		result = "connection"
	}
	return result
}

// storeWindowsCredential stores RDP credentials in Windows Credential Manager
func (l *Launcher) storeWindowsCredential(conn *models.Connection) error {
	if runtime.GOOS != "windows" {
		return nil
	}

	// Build the target name (hostname or hostname:port)
	target := conn.Host
	port := conn.Port
	if port == 0 {
		port = models.ProtocolRDP.GetDefaultPort()
	}
	if port != 3389 {
		target = fmt.Sprintf("%s:%d", conn.Host, port)
	}

	// Build the username (with domain if specified)
	username := conn.Username
	if conn.Domain != "" {
		username = fmt.Sprintf("%s\\%s", conn.Domain, conn.Username)
	}

	// Use cmdkey to store the credential
	// cmdkey /generic:TERMSRV/hostname /user:username /pass:password
	cmd := exec.Command("cmdkey", "/generic:TERMSRV/"+target, "/user:"+username, "/pass:"+conn.Password)
	hideConsoleWindow(cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cmdkey failed: %w, output: %s", err, string(output))
	}

	return nil
}

// RemoveWindowsCredential removes RDP credentials from Windows Credential Manager
func (l *Launcher) RemoveWindowsCredential(conn *models.Connection) error {
	if runtime.GOOS != "windows" {
		return nil
	}

	// Build the target name (hostname or hostname:port)
	target := conn.Host
	port := conn.Port
	if port == 0 {
		port = models.ProtocolRDP.GetDefaultPort()
	}
	if port != 3389 {
		target = fmt.Sprintf("%s:%d", conn.Host, port)
	}

	// Use cmdkey to delete the credential
	// cmdkey /delete:TERMSRV/hostname
	cmd := exec.Command("cmdkey", "/delete:TERMSRV/"+target)
	hideConsoleWindow(cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Don't return error if credential doesn't exist
		if !strings.Contains(string(output), "not found") {
			return fmt.Errorf("cmdkey delete failed: %w, output: %s", err, string(output))
		}
	}

	return nil
}

// CleanupAllCredentials removes all MremoteGO-stored credentials from Windows Credential Manager
func (l *Launcher) CleanupAllCredentials(connections []*models.Connection) error {
	if runtime.GOOS != "windows" {
		return nil
	}

	for _, conn := range connections {
		if conn.Protocol == models.ProtocolRDP {
			_ = l.RemoveWindowsCredential(conn) // Ignore errors, continue with others
		}
	}

	return nil
}
