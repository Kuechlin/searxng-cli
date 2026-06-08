package main

import (
	"os/exec"
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	ti    textinput.Model
	l     list.Model
	focus string
	data  SearchResponse
	err   error
	w     int
	h     int
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Search"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 156

	li := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	li.SetFilteringEnabled(false)
	li.SetShowTitle(false)
	li.SetShowHelp(false)
	li.SetShowStatusBar(false)

	return model{ti: ti, l: li, focus: "input"}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch m.focus {
	case "input":
		m.ti, cmd = m.ti.Update(msg)
		cmds = append(cmds, cmd)
	case "list":
		m.l, cmd = m.l.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height

		m.ti.SetWidth(m.w)
		// set list width/height
		tiHeight := lipgloss.Height(m.ti.View())
		h := m.h - tiHeight - 1
		m.l.SetWidth(m.w)
		m.l.SetHeight(h)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.focus {
			case "input":
				// search
				m.data, m.err = Search(m.ti.Value())
				if m.err != nil {
					break
				}
				var items []list.Item
				for _, v := range m.data.Results {
					items = append(items, v)
				}
				cmd = m.l.SetItems(items)
				cmds = append(cmds, cmd)

				m.focus = "list"
				m.ti.Blur()
				m.l.Select(0)
			case "list":
				// open
				result := m.data.Results[m.l.GlobalIndex()]
				cmd := exec.Command("xdg-open", result.URL)
				if err := cmd.Run(); err != nil {
					panic(err)
				}
			}

		case "s":
			m.focus = "input"
			m.ti.Focus()
			m.l.ResetSelected()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	var c *tea.Cursor
	if !m.ti.VirtualCursor() {
		c = m.ti.Cursor()
	}

	d := lipgloss.NewStyle().Foreground(lipgloss.White)
	line := strings.Repeat("─", max(0, m.w))

	str := lipgloss.JoinVertical(lipgloss.Top, m.ti.View(), d.Render(line), m.l.View())

	v := tea.NewView(str)
	if m.focus == "input" {
		v.Cursor = c
	}
	v.AltScreen = true
	return v
}
