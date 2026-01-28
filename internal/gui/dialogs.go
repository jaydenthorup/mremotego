package gui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jaydenthorup/mremotego/pkg/models"
)

// collectAllFolders recursively collects all folders with their full paths
func (w *MainWindow) collectAllFolders(connections []*models.Connection, prefix string, folderMap map[string]*models.Connection, folderNames *[]string) {
	for _, conn := range connections {
		if conn.IsFolder() {
			var fullPath string
			if prefix == "" {
				fullPath = conn.Name
			} else {
				fullPath = prefix + " / " + conn.Name
			}
			*folderNames = append(*folderNames, fullPath)
			folderMap[fullPath] = conn

			// Recursively process children
			if len(conn.Children) > 0 {
				w.collectAllFolders(conn.Children, fullPath, folderMap, folderNames)
			}
		}
	}
}

// findConnectionParent recursively finds the parent folder and path of a connection
func (w *MainWindow) findConnectionParent(conn *models.Connection, connections []*models.Connection, prefix string) (string, *models.Connection) {
	for _, c := range connections {
		if c.IsFolder() {
			// Check direct children
			for _, child := range c.Children {
				if child == conn {
					var fullPath string
					if prefix == "" {
						fullPath = c.Name
					} else {
						fullPath = prefix + " / " + c.Name
					}
					return fullPath, c
				}
			}

			// Recursively check nested folders
			if len(c.Children) > 0 {
				var fullPath string
				if prefix == "" {
					fullPath = c.Name
				} else {
					fullPath = prefix + " / " + c.Name
				}
				if path, parent := w.findConnectionParent(conn, c.Children, fullPath); parent != nil {
					return path, parent
				}
			}
		}
	}
	return "", nil
}

// findFolderByPath finds a folder by its full path (e.g., "Dev-Ops / Infrastructure / Builders")
func (w *MainWindow) findFolderByPath(path string, connections []*models.Connection) *models.Connection {
	parts := strings.Split(path, " / ")
	current := connections

	for _, part := range parts {
		found := false
		for _, conn := range current {
			if conn.IsFolder() && conn.Name == part {
				if len(parts) == 1 {
					return conn
				}
				current = conn.Children
				parts = parts[1:]
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	return nil
}

// showAddConnectionDialog shows the dialog to add a new connection
func (w *MainWindow) showAddConnectionDialog() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Connection Name")

	protocolSelect := widget.NewSelect([]string{"ssh", "rdp", "vnc", "http", "https", "telnet"}, nil)
	protocolSelect.SetSelected("ssh")

	hostEntry := widget.NewEntry()
	hostEntry.SetPlaceHolder("hostname or IP")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("port (leave empty for default)")

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("username")

	passwordEntry := widget.NewEntry()
	passwordEntry.SetPlaceHolder("password or op://vault/item/field")

	domainEntry := widget.NewEntry()
	domainEntry.SetPlaceHolder("domain (for RDP)")

	descriptionEntry := widget.NewMultiLineEntry()
	descriptionEntry.SetPlaceHolder("Description")

	// Folder selection - recursively collect all folders
	folderNames := []string{"(Root)"}
	folderMap := make(map[string]*models.Connection)
	folderMap["(Root)"] = nil
	w.collectAllFolders(w.manager.GetConfig().Connections, "", folderMap, &folderNames)

	folderSelect := widget.NewSelect(folderNames, nil)
	folderSelect.SetSelected("(Root)")

	// 1Password integration
	storeTo1PasswordCheck := widget.NewCheck("Store password in 1Password", nil)
	vaultSelect := widget.NewSelect([]string{"DevOps", "Private", "Employee"}, nil)
	vaultSelect.SetSelected("DevOps")
	vaultSelect.Hide()

	storeTo1PasswordCheck.OnChanged = func(checked bool) {
		if checked {
			vaultSelect.Show()
		} else {
			vaultSelect.Hide()
		}
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Protocol", Widget: protocolSelect},
			{Text: "Host", Widget: hostEntry},
			{Text: "Port", Widget: portEntry},
			{Text: "Username", Widget: usernameEntry},
			{Text: "Password", Widget: passwordEntry},
			{Text: "Domain", Widget: domainEntry},
			{Text: "Description", Widget: descriptionEntry},
			{Text: "Folder", Widget: folderSelect},
			{Text: "", Widget: storeTo1PasswordCheck},
			{Text: "Vault", Widget: vaultSelect},
		},
		OnSubmit: func() {
			conn := models.NewConnection(nameEntry.Text, models.Protocol(protocolSelect.Selected))
			conn.Host = hostEntry.Text
			conn.Username = usernameEntry.Text
			conn.Password = passwordEntry.Text
			conn.Domain = domainEntry.Text
			conn.Description = descriptionEntry.Text
			conn.Created = time.Now().Format(time.RFC3339)
			conn.Modified = conn.Created

			if portEntry.Text != "" {
				if port, err := strconv.Atoi(portEntry.Text); err == nil {
					conn.Port = port
				}
			} else {
				conn.Port = conn.Protocol.GetDefaultPort()
			}

			// If user wants to store in 1Password, create the item
			if storeTo1PasswordCheck.Checked && conn.Password != "" && !w.manager.IsOnePasswordReference(conn.Password) {
				vault := vaultSelect.Selected
				reference, err := w.manager.CreateOnePasswordItem(vault, conn.Name, conn.Username, conn.Password)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Failed to create 1Password item: %w", err), w.window)
					return
				}
				// Replace password with 1Password reference
				conn.Password = reference
				dialog.ShowInformation("Success", fmt.Sprintf("Password stored in 1Password vault '%s'", vault), w.window)
			}

			// Add to selected folder or root
			selectedFolder := folderSelect.Selected
			if selectedFolder == "(Root)" {
				w.manager.GetConfig().Connections = append(w.manager.GetConfig().Connections, conn)
			} else {
				// Find the folder and add to its children
				folder := folderMap[selectedFolder]
				if folder != nil {
					folder.Children = append(folder.Children, conn)
				}
			}

			if err := w.manager.Save(); err != nil {
				dialog.ShowError(err, w.window)
				return
			}

			w.refreshTree()
			dialog.ShowInformation("Success", "Connection added successfully", w.window)
		},
	}

	d := dialog.NewCustom("Add Connection", "Cancel", form, w.window)
	d.Resize(fyne.NewSize(500, 700))
	d.Show()
}

// showAddFolderDialog shows the dialog to add a new folder
func (w *MainWindow) showAddFolderDialog() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Folder Name")

	// Folder selection - allow creating folders within folders
	folderNames := []string{"(Root)"}
	folderMap := make(map[string]*models.Connection)
	folderMap["(Root)"] = nil
	w.collectAllFolders(w.manager.GetConfig().Connections, "", folderMap, &folderNames)

	folderSelect := widget.NewSelect(folderNames, nil)
	folderSelect.SetSelected("(Root)")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Parent Folder", Widget: folderSelect},
		},
		OnSubmit: func() {
			folder := models.NewFolder(nameEntry.Text)

			// Add to selected folder or root
			selectedFolder := folderSelect.Selected
			if selectedFolder == "(Root)" {
				w.manager.GetConfig().Connections = append(w.manager.GetConfig().Connections, folder)
			} else {
				// Find the folder and add to its children
				parentFolder := folderMap[selectedFolder]
				if parentFolder != nil {
					parentFolder.Children = append(parentFolder.Children, folder)
				}
			}

			if err := w.manager.Save(); err != nil {
				dialog.ShowError(err, w.window)
				return
			}

			w.refreshTree()
			dialog.ShowInformation("Success", "Folder added successfully", w.window)
		},
	}

	dialog.NewCustom("Add Folder", "Cancel", form, w.window).Show()
}

// showEditConnectionDialog shows the dialog to edit a connection
func (w *MainWindow) showEditConnectionDialog(conn *models.Connection) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(conn.Name)

	protocolSelect := widget.NewSelect([]string{"ssh", "rdp", "vnc", "http", "https", "telnet"}, nil)
	protocolSelect.SetSelected(string(conn.Protocol))

	hostEntry := widget.NewEntry()
	hostEntry.SetText(conn.Host)

	portEntry := widget.NewEntry()
	portEntry.SetText(strconv.Itoa(conn.Port))

	usernameEntry := widget.NewEntry()
	usernameEntry.SetText(conn.Username)

	passwordEntry := widget.NewEntry()
	passwordEntry.SetText(conn.Password)

	domainEntry := widget.NewEntry()
	domainEntry.SetText(conn.Domain)

	descriptionEntry := widget.NewMultiLineEntry()
	descriptionEntry.SetText(conn.Description)

	// Folder selection - find current parent folder using recursive search
	currentFolder := "(Root)"
	var parentFolder *models.Connection
	currentFolder, parentFolder = w.findConnectionParent(conn, w.manager.GetConfig().Connections, "")
	if currentFolder == "" {
		currentFolder = "(Root)"
		parentFolder = nil
	}

	folderNames := []string{"(Root)"}
	folderMap := make(map[string]*models.Connection)
	folderMap["(Root)"] = nil
	w.collectAllFolders(w.manager.GetConfig().Connections, "", folderMap, &folderNames)

	folderSelect := widget.NewSelect(folderNames, nil)
	folderSelect.SetSelected(currentFolder)

	// 1Password integration for edit
	storeTo1PasswordCheck := widget.NewCheck("Push password to 1Password", nil)
	vaultSelect := widget.NewSelect([]string{"DevOps", "Private", "Employee"}, nil)
	vaultSelect.SetSelected("DevOps")
	vaultSelect.Hide()

	storeTo1PasswordCheck.OnChanged = func(checked bool) {
		if checked {
			vaultSelect.Show()
		} else {
			vaultSelect.Hide()
		}
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Protocol", Widget: protocolSelect},
			{Text: "Host", Widget: hostEntry},
			{Text: "Port", Widget: portEntry},
			{Text: "Username", Widget: usernameEntry},
			{Text: "Password", Widget: passwordEntry},
			{Text: "Domain", Widget: domainEntry},
			{Text: "Description", Widget: descriptionEntry},
			{Text: "Folder", Widget: folderSelect},
			{Text: "", Widget: storeTo1PasswordCheck},
			{Text: "Vault", Widget: vaultSelect},
		},
		OnSubmit: func() {
			// If user wants to push password to 1Password
			if storeTo1PasswordCheck.Checked && passwordEntry.Text != "" && !w.manager.IsOnePasswordReference(passwordEntry.Text) {
				vault := vaultSelect.Selected
				reference, err := w.manager.CreateOnePasswordItem(vault, nameEntry.Text, usernameEntry.Text, passwordEntry.Text)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Failed to create 1Password item: %w", err), w.window)
					return
				}
				// Replace password with 1Password reference
				conn.Password = reference
				dialog.ShowInformation("Success", fmt.Sprintf("Password stored in 1Password vault '%s'", vault), w.window)
			} else {
				conn.Password = passwordEntry.Text
			}

			conn.Name = nameEntry.Text
			conn.Protocol = models.Protocol(protocolSelect.Selected)
			conn.Host = hostEntry.Text
			conn.Username = usernameEntry.Text
			conn.Domain = domainEntry.Text
			conn.Description = descriptionEntry.Text
			conn.Modified = time.Now().Format(time.RFC3339)

			if port, err := strconv.Atoi(portEntry.Text); err == nil {
				conn.Port = port
			}

			// Handle folder change
			selectedFolder := folderSelect.Selected
			if selectedFolder != currentFolder {
				// Remove from old parent
				if parentFolder != nil {
					// Remove from parent's children
					for i, child := range parentFolder.Children {
						if child == conn {
							parentFolder.Children = append(parentFolder.Children[:i], parentFolder.Children[i+1:]...)
							break
						}
					}
				} else {
					// Remove from root
					for i, c := range w.manager.GetConfig().Connections {
						if c == conn {
							w.manager.GetConfig().Connections = append(w.manager.GetConfig().Connections[:i], w.manager.GetConfig().Connections[i+1:]...)
							break
						}
					}
				}

				// Add to new parent
				if selectedFolder == "(Root)" {
					w.manager.GetConfig().Connections = append(w.manager.GetConfig().Connections, conn)
				} else {
					newParent := folderMap[selectedFolder]
					if newParent != nil {
						newParent.Children = append(newParent.Children, conn)
					}
				}
			}

			if err := w.manager.Save(); err != nil {
				dialog.ShowError(err, w.window)
				return
			}

			w.refreshTree()
			w.updateDetailsPanel(conn)
			dialog.ShowInformation("Success", "Connection updated successfully", w.window)
		},
	}

	d := dialog.NewCustom("Edit Connection", "Cancel", form, w.window)
	d.Resize(fyne.NewSize(500, 700))
	d.Show()
}

// showEditFolderDialog shows the dialog to edit a folder
func (w *MainWindow) showEditFolderDialog(folder *models.Connection) {
	nameEntry := widget.NewEntry()
	nameEntry.SetText(folder.Name)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
		},
		OnSubmit: func() {
			folder.Name = nameEntry.Text

			if err := w.manager.Save(); err != nil {
				dialog.ShowError(err, w.window)
				return
			}

			w.refreshTree()
			dialog.ShowInformation("Success", "Folder updated successfully", w.window)
		},
	}

	dialog.NewCustom("Edit Folder", "Cancel", form, w.window).Show()
}
