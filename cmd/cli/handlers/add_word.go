package handlers

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func addWord() {
	wordInput := widgets.NewParagraph()
	wordInput.Title = "1. Add Word"
	wordInput.SetRect(0, 2, 50, 5)
	wordInput.BorderStyle.Fg = ui.ColorCyan

	instruction := widgets.NewParagraph()
	instruction.Text = "Type your word and press <Enter> to save, or <Esc> to cancel."
	instruction.SetRect(0, 5, 70, 15)
	instruction.Border = false
	instruction.TextStyle = ui.NewStyle(ui.ColorYellow)
	ui.Render(instruction, wordInput)

	inputWord := parseCLIInput(wordInput, instruction)
	if inputWord == EscapeFlag {
		return
	}
	confirm := confirmInput("📤You typed the word: " + inputWord)
	if confirm == EscapeFlag {
		return
	}

}
