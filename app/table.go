package app

import "github.com/charmbracelet/bubbles/table"

func createTable(rows []table.Row, cursor int) table.Model {
	columns := []table.Column{
		{Title: "ETA", Width: 5},
		{Title: "%", Width: 6},
		{Title: "Status", Width: 12},
		{Title: "Down", Width: 8},
		{Title: "Up", Width: 8},
		{Title: "Name", Width: 40},
	}

	t := table.New(
		table.WithColumns(columns),
        table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(6),
	)
    // t.SetCursor(cursor)

	return t
}
