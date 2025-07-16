package handlers

import (
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const EscapeFlag = "QUIT"

func HandleCLISelectedOptions(selectedRow int) {
	ui.Clear()
	switch selectedRow {
	case 0:
		addWord()
	case 1:
	default:
	}
}

func parseCLIInput(placeHolder *widgets.Paragraph, items ...ui.Drawable) string {
	inputWord := ""
	uiEvents := ui.PollEvents()
	items = append(items, placeHolder)
	for {
		ui.Render(items...)
		e := <-uiEvents
		switch e.ID {
		case "<Escape>":
			return EscapeFlag
		case "<Space>":
			inputWord += " "
		case "<Enter>":
			return strings.TrimSpace(inputWord)
		case "<Backspace>":
			if len(inputWord) > 0 {
				inputWord = inputWord[:len(inputWord)-1]
			}
		case "<MouseLeft>", "<MouseRight>", "<MouseMiddle>", "<MouseRelease>":
			continue
		default:
			inputWord += e.ID
		}
		placeHolder.Text = inputWord
	}
}

func confirmInput(input string) string {
	confirm := widgets.NewParagraph()
	confirm.Text = input + "\nPress <Enter> to continue or <Esc> to cancel."
	confirm.SetRect(0, 10, 60, 14)
	confirm.BorderStyle.Fg = ui.ColorGreen
	ui.Clear()
	ui.Render(confirm)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		if e.ID == "<Escape>" {
			return EscapeFlag
		}
		if e.ID == "<Enter>" {
			ui.Clear()
			return ""
		}
	}
}
