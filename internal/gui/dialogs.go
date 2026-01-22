package gui

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/yourusername/mremotego/pkg/models"
)

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

			// Add to root for now (could add folder selection later)
			w.manager.GetConfig().Connections = append(w.manager.GetConfig().Connections, conn)

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

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
		},
		OnSubmit: func() {
			folder := models.NewFolder(nameEntry.Text)

			// Add to root
			w.manager.GetConfig().Connections = append(w.manager.GetConfig().Connections, folder)

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
		},
		OnSubmit: func() {
			conn.Name = nameEntry.Text
			conn.Protocol = models.Protocol(protocolSelect.Selected)
			conn.Host = hostEntry.Text
			conn.Username = usernameEntry.Text
			conn.Password = passwordEntry.Text
			conn.Domain = domainEntry.Text
			conn.Description = descriptionEntry.Text
			conn.Modified = time.Now().Format(time.RFC3339)

			if port, err := strconv.Atoi(portEntry.Text); err == nil {
				conn.Port = port
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
	d.Resize(fyne.NewSize(500, 600))
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
