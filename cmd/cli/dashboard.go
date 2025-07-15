package cli

import (
	"fmt"
	"log"
	"wordrop/cmd/cli/handlers"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func StartDashboard() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	title := setupTitle()
	menu := setupMenu()
	msg := setupMessage()
	uiEvents := ui.PollEvents()

	for {
		ui.Clear()
		ui.Render(title, menu, msg)
		e := <-uiEvents
		switch e.ID {
		case "<Down>":
			menu.ScrollDown()
		case "<Up>":
			menu.ScrollUp()
		case "q", "<C-c>":
			return
		case "<Enter>":
			selectedRowNum := menu.SelectedRow
			if selectedRowNum == len(menu.Rows)-1 {
				return
			}
			msg.Text = fmt.Sprintf("You selected: %s", menu.Rows[selectedRowNum])
			handlers.HandleCLISelectedOptions(selectedRowNum)
		}
	}
}

func setupTitle() *widgets.Paragraph {
	title := widgets.NewParagraph()
	title.Text = "Wordrop Admin Dashboard"
	title.TextStyle = ui.NewStyle(ui.ColorBlue)
	title.SetRect(0, 0, 50, 3)
	title.Border = false
	return title
}

func setupMenu() *widgets.List {
	menu := widgets.NewList()
	menu.Title = "Menu"
	menu.Rows = []string{
		"1. Add Word",
		"2. Add Meaning",
		"3. Add Example",
		"4. Send Daily Word",
		"5. View Word List",
		"q. Quit",
	}
	menu.TitleStyle = ui.NewStyle(ui.ColorGreen)
	menu.SetRect(0, 3, 50, 12)
	menu.TextStyle = ui.NewStyle(ui.ColorBlue, ui.ColorClear, ui.ModifierBold)
	menu.SelectedRowStyle = ui.NewStyle(ui.ColorBlack, ui.ColorWhite)
	return menu
}

func setupMessage() *widgets.Paragraph {
	msg := widgets.NewParagraph()
	msg.Text = "Press <Enter> to select an option."
	msg.SetRect(0, 12, 50, 15)
	msg.Border = false
	msg.TextStyle = ui.NewStyle(ui.ColorGreen)
	return msg
}
