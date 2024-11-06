package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

func GetQuitModel() QuitModel { return QuitModel{} }

// Used to avoid pointer
// If the type of prevModel matches this type, quit the application
type QuitModel struct{}

func (q QuitModel) Init() tea.Cmd                       { return tea.Quit }
func (q QuitModel) Update(tea.Msg) (tea.Model, tea.Cmd) { return q, func() tea.Msg { return tea.Quit } }
func (q QuitModel) View() string                        { return "" }
