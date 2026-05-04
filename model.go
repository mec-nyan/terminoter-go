package main

import (
	"fmt"
	"strings"
	"unicode"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	Options
	Data
	Focused int
	Mode
	Previous rune
	NewNote  []rune
	Size
}

type Mode int

const (
	Normal Mode = iota
	Insert
	Edit
	ConfirmDelete
	Deleted
)

type Size struct {
	Width  int
	Height int
}

func initialModel(opts Options) model {
	return model{
		Options: opts,
		Size:    Size{80, 25},
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		data, err := LoadNotes(m.file)
		return Loader{
			Data:  *data,
			error: err,
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case Loader:
		if msg.error != nil {
			panic(fmt.Sprintf("error loading file: %v", msg.error))
		}

		m.Data = msg.Data
		return m, nil

	case tea.KeyPressMsg:

		switch m.Mode {

		case Normal:

			switch msg.Code {

			case '\x1b', 'q':
				return m, tea.Quit

			case 'j':
				m.Focused++
				if m.Focused == len(m.Notes) {
					m.Focused = 0
				}

			case 'k':
				m.Focused--
				if m.Focused < 0 {
					m.Focused = len(m.Notes) - 1
				}

			case 'i':
				m.Mode = Insert

			case 'd':
				if len(m.Notes) > 0 {
					m.Mode = ConfirmDelete
				}
			}

		case Insert:
			switch msg.String() {
			case "enter":
				m.Notes = append(m.Notes, Note{Content: string(m.NewNote)})
				m.NewNote = nil
				m.Mode = Normal
				// Focus the last added item.
				m.Focused = len(m.Notes) - 1

			case "esc":
				m.NewNote = nil
				m.Mode = Normal

			// Basic line editing.
			case "backspace":
				m.NewNote = m.NewNote[:len(m.NewNote)-1]

			case "ctrl+u":
				m.NewNote = nil

			case "ctrl+w":
				for len(m.NewNote) > 0 {
					m.NewNote = m.NewNote[:len(m.NewNote)-1]
					if len(m.NewNote) == 0 {
						break
					}

					last := m.NewNote[len(m.NewNote)-1]
					if last == ' ' {
						break
					}
				}

			default:
				if unicode.IsGraphic(msg.Code) {
					m.NewNote = append(m.NewNote, []rune(msg.Text)...)
				}
			}

		case ConfirmDelete:
			switch msg.String() {
			case "enter", "y":
				// NOTE: Move this logic to a method (i.e. Notes.delete(n)).
				if len(m.Notes) == 1 {
					m.Notes = nil
				} else if m.Focused == 0 {
					m.Notes = m.Notes[1:]
				} else if m.Focused == len(m.Notes)-1 {
					m.Notes = m.Notes[:len(m.Notes)-1]
				} else {
					m.Notes = append(m.Notes[:m.Focused], m.Notes[m.Focused+1:]...)
				}
				if m.Focused > 0 {
					m.Focused--
				}
				m.Mode = Deleted

			case "esc":
				m.Mode = Normal

			default:
				m.Mode = Normal
			}

		case Deleted:
			m.Mode = Normal

		default:
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Yellow).
		Width(m.Width).
		Align(lipgloss.Center).
		MarginBottom(1)

	noteStyle := lipgloss.NewStyle().
		Width(m.Width-2).
		Foreground(lipgloss.Blue).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Blue).
		Padding(0, 2).
		Margin(1, 1, 0, 1)

	focusedNoteStyle := noteStyle.
		BorderForeground(lipgloss.Green).
		Foreground(lipgloss.Green)

	deleteStyle := noteStyle.
		BorderForeground(lipgloss.Red).
		Foreground(lipgloss.Red).
		Margin(1)

	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Red).
		PaddingTop(1).
		PaddingLeft(2)

	var s strings.Builder
	s.WriteString(titleStyle.Render("  Notes"))

	switch m.Mode {
	case Normal:
		for i, note := range m.Notes {
			if i == m.Focused {
				s.WriteString(focusedNoteStyle.Render(note.Content))
			} else {
				s.WriteString(noteStyle.Render(note.Content))
			}
		}

	case Insert:
		fmt.Fprintf(&s, "\nNew note:\n\n > %s█", string(m.NewNote))

	case ConfirmDelete:
		s.WriteString(warningStyle.Render("Delete this note?\nThis action can't be undone."))
		s.WriteString(deleteStyle.Render(m.Notes[m.Focused].Content))
		s.WriteString(warningStyle.Render("Press y or <enter> to confirm. <esc> to cancel."))

	case Deleted:
		s.WriteString(deleteStyle.Render("Note deleted. Press any key."))

	default:
	}

	v := tea.NewView(s.String())

	v.AltScreen = true

	return v
}
