package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/yourusername/mremotego/pkg/models"
)

// createContextMenu creates a context menu for the selected connection
func (w *MainWindow) createContextMenu(conn interface{}) *fyne.Container {
	if conn == nil {
		return container.NewVBox()
	}

	c, ok := conn.(*models.Connection)
	if !ok {
		return container.NewVBox()
	}

	if c.IsFolder() {
		return container.NewVBox(
			widget.NewButton("Edit Folder", func() {
				w.showEditFolderDialog(c)
			}),
			widget.NewButton("Delete Folder", func() {
				w.deleteSelected()
			}),
		)
	}

	return container.NewVBox(
		widget.NewButton("Connect", func() {
			w.connectToConnection(c)
		}),
		widget.NewSeparator(),
		widget.NewButton("Edit", func() {
			w.showEditConnectionDialog(c)
		}),
		widget.NewButton("Duplicate", func() {
			// TODO: Implement duplicate
		}),
		widget.NewSeparator(),
		widget.NewButton("Delete", func() {
			w.deleteSelected()
		}),
	)
}
