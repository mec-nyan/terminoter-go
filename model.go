package main

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// model is the state of our app.
type model struct {
	Options
	Data
	Focused int
	Mode
	Size
	// This will be in their own struct in a future.
	NewNote []rune
	ConfirmInsert
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

type ConfirmInsert struct {
	First, Second bool
}

func (c *ConfirmInsert) Set() {
	if !c.First {
		c.First = true
	} else if !c.Second {
		c.Second = true
	}
}

func (c *ConfirmInsert) UnSet() {
	c.First = false
	c.Second = false
}

func (c *ConfirmInsert) Confirm() bool {
	return c.First && c.Second
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
				if m.Confirm() {
					content := strings.TrimSpace(string(m.NewNote))
					m.Notes = append(m.Notes, Note{Content: content})
					m.NewNote = nil
					m.Mode = Normal
					// Focus the last added item.
					m.Focused = len(m.Notes) - 1
					m.UnSet()
					// TODO: Better error handling.
					err := SaveNotes(&m.Data, m.file)
					if err != nil {
						log.Fatalf("Error saving notes: %v", err)
					}
				} else {
					// NOTE: Maybe a more descriptive name here...
					m.Set()
					// We trim the extra whitespace on note insertion.
					m.NewNote = append(m.NewNote, '\n')
				}

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
					m.UnSet()
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
				m.Mode = Normal

			case "esc":
				m.Mode = Normal

			default:
				m.Mode = Normal
			}

		default:
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Cyan).
		Width(m.Width).
		Align(lipgloss.Center).
		MarginBottom(1)

	noteStyle := lipgloss.NewStyle().
		Width(m.Width-4).
		Foreground(lipgloss.Blue).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Blue).
		Padding(0, 2).
		Margin(0, 2)

	focusedNoteStyle := noteStyle.
		BorderForeground(lipgloss.Green).
		Foreground(lipgloss.Green)

	newNoteStyle := noteStyle.
		BorderForeground(lipgloss.Cyan).
		Foreground(lipgloss.White)

	deleteStyle := noteStyle.
		BorderForeground(lipgloss.Red).
		Foreground(lipgloss.White)

	warningStyle := lipgloss.NewStyle().
		Width(m.Width-4).
		Foreground(lipgloss.Red).
		Italic(true).
		Border(lipgloss.BlockBorder(), false, false, false, true).
		BorderForeground(lipgloss.Red).
		Padding(0, 2).
		Margin(1, 2)

	title := titleStyle.Render("== Notes ==")

	// Show title in default style.
	notes := []string{title}

	// Add the notes. Highlight the focused notes according to current mode.
	var showStyle lipgloss.Style

	for i, note := range m.Notes {
		if i == m.Focused {
			switch m.Mode {
			case ConfirmDelete:
				// Highilght the note about to be deleted (or not).
				showStyle = deleteStyle
				notes = append(notes, warningStyle.Render(
					"  Delete this note? You want be able to recover it later.\n\n"+
						"Hit <Escape> to cancel. <Enter> (or 'y') to confirm deletion."))
			case Insert:
				// When in Insert mode, the focus is on the note being inserted.
				showStyle = noteStyle
			default:
				showStyle = focusedNoteStyle
			}
		} else {
			showStyle = noteStyle
		}

		notes = append(notes, showStyle.Render(note.Content))
	}

	if m.Mode == Insert {
		// Add the new input field where the new note will be.
		notes = append(notes, newNoteStyle.Render(fmt.Sprintf(
			"New note:\n\n%s█", string(m.NewNote))))

	}

	// TODO: Show a message at the bottom when a note has been deleted.

	v := tea.NewView(lipgloss.JoinVertical(
		lipgloss.Left, notes...))

	v.AltScreen = true

	return v
}
