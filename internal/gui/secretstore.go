package gui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/jaydenthorup/mremotego/internal/secrets"
)

// passwordPlaceholder documents every accepted form of the password field.
const passwordPlaceholder = "password, op://vault/item/field or bw://item-id"

// secretStoreControls are the shared widgets that let a connection dialog pick
// an existing secret or push a new one to a password manager. Both the add and
// the edit dialog use them, so the two stay in step.
type secretStoreControls struct {
	window *MainWindow

	opCheck     *widget.Check
	bwCheck     *widget.Check
	vaultSelect *widget.Select
	pickButton  *widget.Button
}

// newSecretStoreControls builds the controls for one dialog. usernameEntry and
// passwordEntry are filled in when the user picks an existing item.
func (w *MainWindow) newSecretStoreControls(usernameEntry, passwordEntry *widget.Entry) *secretStoreControls {
	c := &secretStoreControls{window: w}

	c.vaultSelect = widget.NewSelect([]string{"DevOps", "Private", "Employee"}, nil)
	c.vaultSelect.SetSelected("DevOps")
	c.vaultSelect.Hide()

	c.opCheck = widget.NewCheck("Store password in 1Password", nil)
	c.bwCheck = widget.NewCheck("Store password in Bitwarden", nil)

	// Storing the same password in two vaults would leave one of them stale, so
	// the two options exclude each other.
	c.opCheck.OnChanged = func(checked bool) {
		if checked {
			c.bwCheck.SetChecked(false)
			c.vaultSelect.Show()
		} else {
			c.vaultSelect.Hide()
		}
	}

	bitwarden := w.bitwardenProvider()

	c.bwCheck.OnChanged = func(checked bool) {
		if !checked {
			return
		}
		c.opCheck.SetChecked(false)
		// Starting the helper process takes a moment; do it now so that saving
		// the form does not block.
		if bitwarden != nil {
			bitwarden.Warmup()
		}
	}

	c.pickButton = widget.NewButton("Bitwarden...", func() {
		w.showBitwardenPicker(func(item secrets.BitwardenItem) {
			passwordEntry.SetText("bw://" + item.ID)
			if usernameEntry.Text == "" {
				usernameEntry.SetText(item.Username)
			}
		})
	})

	if bitwarden == nil || !bitwarden.IsEnabled() {
		c.bwCheck.SetText("Store password in Bitwarden (bw CLI not installed)")
		c.bwCheck.Disable()
		c.pickButton.Hide()
	}

	return c
}

// passwordWidget wraps the password entry so the picker button sits next to it.
func (c *secretStoreControls) passwordWidget(passwordEntry *widget.Entry) fyne.CanvasObject {
	return container.NewBorder(nil, nil, nil, c.pickButton, passwordEntry)
}

// formItems returns the rows to append to a connection form.
func (c *secretStoreControls) formItems() []*widget.FormItem {
	return []*widget.FormItem{
		{Text: "", Widget: c.opCheck},
		{Text: "", Widget: c.bwCheck},
		{Text: "Vault", Widget: c.vaultSelect},
	}
}

// storeIfRequested pushes the password to the selected password manager and
// returns the value to write to the configuration: either a reference to the
// newly created item, or the password unchanged. providerName names the manager
// that was used, or is empty when nothing was stored.
func (c *secretStoreControls) storeIfRequested(title, username, password, protocol, host string) (value string, providerName string, err error) {
	// Nothing to store, or the value already is a reference.
	if password == "" || c.window.manager.IsSecretReference(password) {
		return password, "", nil
	}

	req := secrets.CreateItemRequest{
		Title:    title,
		Username: username,
		Password: password,
	}

	var scheme string
	switch {
	case c.opCheck.Checked:
		scheme, providerName = secrets.SchemeOnePassword, "1Password"
		req.Vault = c.vaultSelect.Selected
	case c.bwCheck.Checked:
		scheme, providerName = secrets.SchemeBitwarden, "Bitwarden"
		// Recording the target makes the item useful in the Bitwarden clients
		// as well, not just here.
		if protocol != "" && host != "" {
			req.URI = protocol + "://" + host
		}
	default:
		return password, "", nil
	}

	reference, err := c.window.manager.CreateSecretItem(scheme, req)
	if err != nil {
		return "", providerName, fmt.Errorf("failed to create %s item: %w", providerName, err)
	}

	return reference, providerName, nil
}

// bitwardenProvider returns the Bitwarden provider, or nil when it is not
// registered.
func (w *MainWindow) bitwardenProvider() *secrets.BitwardenProvider {
	provider, ok := w.launcher.Secrets().ByScheme(secrets.SchemeBitwarden)
	if !ok {
		return nil
	}

	bitwarden, _ := provider.(*secrets.BitwardenProvider)
	return bitwarden
}

// showBitwardenPicker lets the user choose a login item from their vault. The
// chosen item is reported as a reference, so the password itself never reaches
// the configuration file.
func (w *MainWindow) showBitwardenPicker(onPick func(secrets.BitwardenItem)) {
	bitwarden := w.bitwardenProvider()
	if bitwarden == nil || !bitwarden.IsEnabled() {
		dialog.ShowError(fmt.Errorf("Bitwarden CLI (bw) is not installed or not in PATH"), w.window)
		return
	}

	var all []secrets.BitwardenItem
	var filtered []secrets.BitwardenItem
	var selected *secrets.BitwardenItem

	statusLabel := widget.NewLabel("Loading vault items...")

	list := widget.NewList(
		func() int { return len(filtered) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(filtered) {
				return
			}
			item := filtered[id]
			label := item.Name
			if item.Username != "" {
				label += "  -  " + item.Username
			}
			obj.(*widget.Label).SetText(label)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		if id < len(filtered) {
			item := filtered[id]
			selected = &item
		}
	}

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search by name, username or URL")

	// The whole vault is fetched once and filtered locally: asking the CLI on
	// every keystroke would mean a round trip per character.
	applyFilter := func(query string) {
		query = strings.ToLower(strings.TrimSpace(query))

		filtered = filtered[:0]
		for _, item := range all {
			if query == "" ||
				strings.Contains(strings.ToLower(item.Name), query) ||
				strings.Contains(strings.ToLower(item.Username), query) ||
				strings.Contains(strings.ToLower(item.URI), query) {
				filtered = append(filtered, item)
			}
		}

		selected = nil
		list.UnselectAll()
		list.Refresh()
		statusLabel.SetText(fmt.Sprintf("%d of %d items", len(filtered), len(all)))
	}

	searchEntry.OnChanged = applyFilter

	// load fetches the vault in the background; every UI update goes through
	// fyne.Do because Fyne widgets may only be touched from the UI goroutine.
	load := func(sync bool) {
		statusLabel.SetText("Loading vault items...")

		go func() {
			if sync {
				if err := bitwarden.Sync(); err != nil {
					fyne.Do(func() {
						statusLabel.SetText("Sync failed: " + err.Error())
					})
					return
				}
			}

			items, err := bitwarden.ListLoginItems("")
			fyne.Do(func() {
				if err != nil {
					statusLabel.SetText("Failed to load items")
					dialog.ShowError(err, w.window)
					return
				}

				all = items
				applyFilter(searchEntry.Text)
			})
		}()
	}

	syncButton := widget.NewButton("Sync vault", func() { load(true) })

	content := container.NewBorder(
		searchEntry,
		container.NewHBox(statusLabel, layout.NewSpacer(), syncButton),
		nil,
		nil,
		list,
	)

	picker := dialog.NewCustomConfirm("Pick from Bitwarden", "Use", "Cancel", content,
		func(confirmed bool) {
			if confirmed && selected != nil {
				onPick(*selected)
			}
		}, w.window)
	picker.Resize(fyne.NewSize(600, 500))
	picker.Show()

	load(false)
}
