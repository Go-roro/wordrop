package cli

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"strings"
)

func HandleCLISelectedOptions(selectedRow int) {
	ui.Clear()
	switch selectedRow {
	case 0:
		addWord()
	case 1:
	default:
	}
}

func addWord() {
	wordInput := widgets.NewParagraph()
	wordInput.Title = "1. Add Word"
	wordInput.SetRect(0, 3, 50, 6)
	wordInput.BorderStyle.Fg = ui.ColorCyan

	instruction := widgets.NewParagraph()
	instruction.Text = "Type your word and press <Enter> to save, or <Esc> to cancel."
	instruction.SetRect(0, 5, 70, 15)
	instruction.Border = false
	instruction.TextStyle = ui.NewStyle(ui.ColorYellow)
	ui.Render(instruction, wordInput)

	inputWord := ""
	uiEvents := ui.PollEvents()
	for {
		ui.Render(instruction, wordInput)
		e := <-uiEvents
		switch e.ID {
		case "<Escape>":
			return
		case "<Space>":
			inputWord += " "
		case "<Enter>":
			finalInput := strings.TrimSpace(inputWord)
			if inputWord != "" {
				showAddConfirmation(finalInput)
				return
			}
		case "<Backspace>":
			if len(inputWord) > 0 {
				inputWord = inputWord[:len(inputWord)-1]
			}
		case "<MouseLeft>", "<MouseRight>", "<MouseMiddle>", "<MouseRelease>":
			continue
		default:
			inputWord += e.ID
		}
		wordInput.Text = inputWord
	}
}

func showAddConfirmation(input string) {
	confirm := widgets.NewParagraph()
	confirm.Text = "✅Word added: " + input + "\nPress <Enter> to continue or <Esc> to cancel."
	confirm.SetRect(0, 10, 60, 14)
	confirm.BorderStyle.Fg = ui.ColorGreen
	ui.Clear()
	ui.Render(confirm)

	uiEvents := ui.PollEvents()
	for e := range uiEvents {
		if e.ID == "<Enter>" || e.ID == "<Escape>" {
			return
		}
	}
}
