package main

import (
	"os/exec"
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	ti      textinput.Model
	l       list.Model
	s       spinner.Model
	focus   string
	loading bool
	data    SearchResponse
	err     error
	w       int
	h       int
}

type fetchDoneMsg struct {
	data SearchResponse
	err  error
}

func startFetchCmd(promt string) tea.Cmd {
	return func() tea.Msg {
		data, err := Search(promt)
		return fetchDoneMsg{data: data, err: err}
	}
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

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Magenta)

	return model{ti: ti, l: li, s: s, focus: "input"}
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
				m.loading = true
				m.ti.Prompt = ""
				cmds = append(cmds, startFetchCmd(m.ti.Value()), m.s.Tick)

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
	case fetchDoneMsg:
		m.loading = false
		m.ti.Prompt = "> "
		if m.err != nil {
			m.err = msg.err
			break
		}
		m.data = msg.data
		var items []list.Item
		for _, v := range m.data.Results {
			items = append(items, v)
		}
		cmd = m.l.SetItems(items)
		cmds = append(cmds, cmd)
	}

	if m.loading {
		m.s, cmd = m.s.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	if m.err != nil {
		return tea.NewView(m.err.Error())
	}

	var c *tea.Cursor
	if !m.ti.VirtualCursor() {
		c = m.ti.Cursor()
	}

	d := lipgloss.NewStyle().Foreground(lipgloss.White)
	line := strings.Repeat("─", max(0, m.w))

	var header string
	if m.loading {
		header = lipgloss.JoinHorizontal(lipgloss.Left, m.s.View(), m.ti.View())
	} else {
		header = m.ti.View()
	}

	str := lipgloss.JoinVertical(lipgloss.Top, header, d.Render(line), m.l.View())

	v := tea.NewView(str)
	if m.focus == "input" {
		v.Cursor = c
	}
	v.AltScreen = true
	return v
}
