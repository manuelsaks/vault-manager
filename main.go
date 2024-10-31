package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Struktura przechowująca style
type Styles struct {
	HeaderTextColor tcell.Color
	HeaderAlign     int
	CellTextColor   tcell.Color
	CellAlign       int
	Selectable      bool
	Expansion       int
}

// Funkcja do utworzenia nowej struktury z domyślnymi stylami
func DefaultStyles() Styles {
	return Styles{
		HeaderTextColor: tcell.ColorWhite,
		HeaderAlign:     tview.AlignLeft,
		CellTextColor:   tcell.ColorLightBlue,
		CellAlign:       tview.AlignLeft,
		Selectable:      true,
		Expansion:       1,
	}
}

func createHeader(styles Styles) *tview.TextView {
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetText("SECRETS[2]").
		SetTextAlign(tview.AlignCenter)
	return header
}

func createTable(styles Styles, headers []string, data [][]string) *tview.Table {
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)

	for i, header := range headers {
		table.SetCell(0, i, tview.NewTableCell(header).
			SetSelectable(false).
			SetAlign(styles.HeaderAlign).
			SetExpansion(styles.Expansion).
			SetTextColor(styles.HeaderTextColor))
	}

	for i, row := range data {
		for j, cell := range row {
			table.SetCell(i+1, j, tview.NewTableCell(cell).
				SetAlign(styles.CellAlign).
				SetExpansion(styles.Expansion).
				SetTextColor(styles.CellTextColor).
				SetSelectable(styles.Selectable))
		}
	}

	return table
}

func createSearchBar() *tview.InputField {
	searchBar := tview.NewInputField()
	searchBar.SetLabel("|>")
	searchBar.SetFieldWidth(0)
	return searchBar
}

func createVersionTable(styles Styles, headers []string, currentVersion []string, olderVersions [][]string) *tview.Table {
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false)

	versionHeaders := []string{"VERSION", "STATUS", "ACTIVATION_DATE", "EXPIRATION_DATE", "CREATED", "UPDATED"}
	for i, header := range versionHeaders {
		table.SetCell(0, i, tview.NewTableCell(header).
			SetSelectable(false).
			SetAlign(styles.HeaderAlign).
			SetExpansion(styles.Expansion).
			SetTextColor(styles.HeaderTextColor))
	}

	rowIndex := 1
	table.SetCell(rowIndex, 0, tview.NewTableCell("CURRENT VERSION").
		SetSelectable(false).
		SetAlign(styles.HeaderAlign).
		SetExpansion(styles.Expansion).
		SetTextColor(styles.HeaderTextColor))
	rowIndex++
	for i, cell := range currentVersion {
		table.SetCell(rowIndex, i, tview.NewTableCell(cell).
			SetAlign(styles.CellAlign).
			SetExpansion(styles.Expansion).
			SetTextColor(styles.CellTextColor).
			SetSelectable(styles.Selectable))
	}

	rowIndex += 2
	table.SetCell(rowIndex, 0, tview.NewTableCell("OLDER VERSIONS").
		SetSelectable(false).
		SetAlign(styles.HeaderAlign).
		SetExpansion(styles.Expansion).
		SetTextColor(styles.HeaderTextColor))
	rowIndex++
	for _, row := range olderVersions {
		for i, cell := range row {
			table.SetCell(rowIndex, i, tview.NewTableCell(cell).
				SetAlign(styles.CellAlign).
				SetExpansion(styles.Expansion).
				SetTextColor(styles.CellTextColor).
				SetSelectable(styles.Selectable))
		}
		rowIndex++
	}

	return table
}

func updateTable(header *tview.TextView, table *tview.Table, headers []string, data [][]string, styles Styles, filter string) {
	table.Clear()
	for i, headerText := range headers {
		table.SetCell(0, i, tview.NewTableCell(headerText).
			SetSelectable(false).
			SetAlign(styles.HeaderAlign).
			SetExpansion(styles.Expansion).
			SetTextColor(styles.HeaderTextColor))
	}
	rowIndex := 1
	filter = strings.ToLower(filter)
	for _, row := range data {
		if filter == "" || strings.Contains(strings.ToLower(row[0]), filter) {
			for i, cell := range row {
				table.SetCell(rowIndex, i, tview.NewTableCell(cell).
					SetAlign(styles.CellAlign).
					SetExpansion(styles.Expansion).
					SetTextColor(styles.CellTextColor).
					SetSelectable(styles.Selectable))
			}
			rowIndex++
		}
	}
	header.SetText(fmt.Sprintf("SECRETS[%d]_<%s>", rowIndex-1, filter))
}

func handleInput(app *tview.Application, flex *tview.Flex, searchBar *tview.InputField, header *tview.TextView, table *tview.Table, headers []string, data [][]string, styles Styles, isSearchBarVisible *bool) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == '/' {
			if !*isSearchBarVisible {
				flex.Clear()
				flex.AddItem(searchBar, 1, 1, true).
					AddItem(header, 1, 1, false).
					AddItem(table, 0, 10, true)
				*isSearchBarVisible = true
				app.SetFocus(searchBar)
			} else {
				flex.Clear()
				flex.AddItem(header, 1, 1, false).
					AddItem(table, 0, 10, true)
				*isSearchBarVisible = false
				app.SetFocus(table)
			}
			return nil
		}
		if event.Key() == tcell.KeyEscape {
			searchBar.SetText("")
			updateTable(header, table, headers, data, styles, "")
		}
		return event
	})
}

func main() {
	app := tview.NewApplication()
	styles := DefaultStyles()
	headers := []string{"NAME", "TYPE", "STATUS", "EXPIRES"}
	data := [][]string{
		{"test1", "text", "enabled", "2024-01-02"},
		{"test2", "null", "enabled", "null"},
		{"dev", "null", "enabled", "null"},
	}

	header := createHeader(styles)
	table := createTable(styles, headers, data)
	searchBar := createSearchBar()
	isSearchBarVisible := false

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(header, 1, 1, false).
		AddItem(table, 0, 10, true)

	handleInput(app, flex, searchBar, header, table, headers, data, styles, &isSearchBarVisible)

	searchBar.SetChangedFunc(func(text string) {
		updateTable(header, table, headers, data, styles, text)
	})
	searchBar.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			updateTable(header, table, headers, data, styles, searchBar.GetText())
			flex.Clear()
			flex.AddItem(header, 1, 1, false).
				AddItem(table, 0, 10, true)
			isSearchBarVisible = false
			app.SetFocus(table)
		} else if key == tcell.KeyEscape {
			searchBar.SetText("")
			updateTable(header, table, headers, data, styles, "")
			flex.Clear()
			flex.AddItem(header, 1, 1, false).
				AddItem(table, 0, 10, true)
			isSearchBarVisible = false
			app.SetFocus(table)
		}
	})

	table.SetSelectedFunc(func(row, column int) {
		if row > 0 {
			// Example data for versions
			currentVersion := []string{"cddcssd", "Enabled", "null", "null", "10/30/2024, 2:48:08 PM", "10/30/2024, 2:48:08 PM"}
			olderVersions := [][]string{
				{"xxxxadsx", "Enabled", "null", "null", "10/30/2024, 2:48:08 PM", "10/30/2024, 2:48:08 PM"},
				{"sdfdscdsc", "Enabled", "null", "null", "10/30/2024, 2:48:08 PM", "10/30/2024, 2:48:08 PM"},
				{"fadcdcsa", "Enabled", "null", "null", "10/30/2024, 2:48:08 PM", "10/30/2024, 2:48:08 PM"},
			}

			versionTable := createVersionTable(styles, headers, currentVersion, olderVersions)
			flex.Clear()
			flex.AddItem(header, 1, 1, false).
				AddItem(versionTable, 0, 10, true)
			app.SetFocus(versionTable)
		}
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
