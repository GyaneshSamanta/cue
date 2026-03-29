package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles shared across all TUI components.
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF6B6B")).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C0C0C0"))

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7DCFFF")).
			Bold(true)

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A9B1D6"))

	DimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#565F89"))

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9ECE6A"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F7768E"))

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E0AF68"))

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7AA2F7")).
			Padding(1, 2)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#BB9AF7")).
			MarginBottom(1)
)

// PickerItem is an item displayed in a TUI list picker.
type PickerItem struct {
	Name        string
	Description string
	Tag         string // e.g. "[NEW]", "[⚠]"
	SizeMB      int
}

func (i PickerItem) Title() string {
	tag := ""
	if i.Tag != "" {
		tag = WarningStyle.Render(i.Tag) + " "
	}
	return tag + SelectedStyle.Render(i.Name)
}
func (i PickerItem) FilterValue() string { return i.Name }
func (i PickerItem) Description2() string {
	if i.SizeMB > 0 {
		return fmt.Sprintf("%s  (%s)", i.Description, formatSize(i.SizeMB))
	}
	return i.Description
}

func formatSize(mb int) string {
	if mb >= 1000 {
		return fmt.Sprintf("~%.1f GB", float64(mb)/1000)
	}
	return fmt.Sprintf("~%d MB", mb)
}

// --- Simple List Picker ---

type pickerModel struct {
	items    []PickerItem
	cursor   int
	selected int
	search   string
	filtered []int
	title    string
	quitting bool
}

func NewPicker(title string, items []PickerItem) *pickerModel {
	filtered := make([]int, len(items))
	for i := range items {
		filtered[i] = i
	}
	return &pickerModel{
		items:    items,
		filtered: filtered,
		selected: -1,
		title:    title,
	}
}

func (m pickerModel) Init() tea.Cmd { return nil }

func (m *pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.filtered) > 0 {
				m.selected = m.filtered[m.cursor]
			}
			return m, tea.Quit
		case "backspace":
			if len(m.search) > 0 {
				m.search = m.search[:len(m.search)-1]
				m.applyFilter()
			}
		default:
			if len(msg.String()) == 1 {
				m.search += msg.String()
				m.applyFilter()
			}
		}
	}
	return m, nil
}

func (m *pickerModel) applyFilter() {
	m.filtered = m.filtered[:0]
	for i, item := range m.items {
		if m.search == "" || strings.Contains(strings.ToLower(item.Name), strings.ToLower(m.search)) ||
			strings.Contains(strings.ToLower(item.Description), strings.ToLower(m.search)) {
			m.filtered = append(m.filtered, i)
		}
	}
	m.cursor = 0
}

func (m pickerModel) View() string {
	var b strings.Builder

	b.WriteString(HeaderStyle.Render("┌─ "+m.title+" ─────────────────────────────────────────┐") + "\n")
	b.WriteString(DimStyle.Render("  ↑/↓ navigate   Enter select   / search   q quit") + "\n")

	if m.search != "" {
		b.WriteString(fmt.Sprintf("  Search: %s\n", SelectedStyle.Render(m.search)))
	}
	b.WriteString("\n")

	for i, idx := range m.filtered {
		item := m.items[idx]
		cursor := "  "
		nameStyle := NormalStyle
		if i == m.cursor {
			cursor = SelectedStyle.Render("▶ ")
			nameStyle = SelectedStyle
		}

		tag := ""
		if item.Tag != "" {
			tag = WarningStyle.Render(" "+item.Tag)
		}

		size := ""
		if item.SizeMB > 0 {
			size = DimStyle.Render("  " + formatSize(item.SizeMB))
		}

		b.WriteString(fmt.Sprintf("%s%s%s  %s%s\n",
			cursor,
			nameStyle.Render(item.Name),
			tag,
			DimStyle.Render(item.Description),
			size))
	}

	if len(m.filtered) == 0 {
		b.WriteString(DimStyle.Render("  No matches found.\n"))
	}

	b.WriteString(HeaderStyle.Render("└─────────────────────────────────────────────────────────┘"))
	return b.String()
}

// Selected returns the selected item index, or -1 if cancelled.
func (m pickerModel) Selected() int { return m.selected }

// RunPicker displays an interactive picker and returns the selected index.
func RunPicker(title string, items []PickerItem) (int, error) {
	p := NewPicker(title, items)
	model, err := tea.NewProgram(p).Run()
	if err != nil {
		return -1, err
	}
	return model.(*pickerModel).Selected(), nil
}

// --- Confirm Prompt ---

type confirmModel struct {
	question string
	yes      bool
	done     bool
}

func (m confirmModel) Init() tea.Cmd { return nil }

func (m *confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y", "enter":
			m.yes = true
			m.done = true
			return m, tea.Quit
		case "n", "N", "q", "ctrl+c":
			m.yes = false
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	return fmt.Sprintf("%s %s",
		SelectedStyle.Render(m.question),
		DimStyle.Render("[Y/n]"))
}

// RunConfirm displays a yes/no prompt.
func RunConfirm(question string) (bool, error) {
	m := &confirmModel{question: question}
	model, err := tea.NewProgram(m).Run()
	if err != nil {
		return false, err
	}
	return model.(*confirmModel).yes, nil
}

// --- Text Input ---

// RunTextInput shows a text input prompt and returns the entered value.
func RunTextInput(prompt string, placeholder string) (string, error) {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()

	m := &textInputModel{input: ti, prompt: prompt}
	model, err := tea.NewProgram(m).Run()
	if err != nil {
		return "", err
	}
	return model.(*textInputModel).value, nil
}

type textInputModel struct {
	input  textinput.Model
	prompt string
	value  string
	done   bool
}

func (m textInputModel) Init() tea.Cmd { return textinput.Blink }

func (m *textInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.value = m.input.Value()
			m.done = true
			return m, tea.Quit
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m textInputModel) View() string {
	return fmt.Sprintf("%s\n%s", HeaderStyle.Render(m.prompt), m.input.View())
}

// Suppress unused import warnings
var _ = list.NewDefaultDelegate
