package main

import (
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	inputField := tview.NewInputField().
		SetLabel("> ").
		SetPlaceholder("Search: ").
		SetFieldWidth(0)

	resultsList := tview.NewList()

	errorText := tview.NewTextView().SetTextColor(tcell.ColorRed)

	inputField.SetDoneFunc(func(key tcell.Key) {
		data, err := Search(inputField.GetText())
		if err != nil {
			errorText.SetText(err.Error())
			return
		}
		errorText.SetText("no error")

		resultsList.Clear()
		for _, result := range data.Results {
			resultsList.AddItem(result.Title, result.Content, 0, func() {
				cmd := exec.Command("xdg-open", result.URL)
				if err := cmd.Run(); err != nil {
					panic(err)
				}
			})
		}
		app.SetFocus(resultsList)
	})

	grid := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputField, inputField.GetFieldHeight(), 0, true).
		AddItem(errorText, errorText.GetFieldHeight(), 0, false).
		AddItem(resultsList, 0, 1, false)

	if err := app.SetRoot(grid, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}
