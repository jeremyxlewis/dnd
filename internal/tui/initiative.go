package tui

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// combatant represents a participant in combat.
type combatant struct {
	name       string
	init       int
	hp         int
	maxHP      int
	conditions []string
}

// initiativeTracker manages combat initiative and HP/conditions.
type initiativeTracker struct {
	list       list.Model
	combatants []combatant
	width      int
	height     int
	inputMode  string // "add_name", "add_init", "hp_damage", etc.
	textInput  textinput.Model
	selected   int
}

func newInitiativeTracker(width, height int) initiativeTracker {
	ti := textinput.New()
	ti.Placeholder = "Enter name"
	ti.Focus()

	items := []list.Item{}
	l := list.New(items, customDelegate{}, width, height-ListHeightPadding)
	l.Title = "Combatants (sorted by initiative)"

	return initiativeTracker{
		list:       l,
		combatants: []combatant{},
		width:      width,
		height:     height,
		inputMode:  InputModeAddName,
		textInput:  ti,
		selected:   0,
	}
}

func (m initiativeTracker) Init() tea.Cmd {
	return textinput.Blink
}

func (m initiativeTracker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.inputMode == InputModeAddName {
				name := m.textInput.Value()
				if name != "" {
					m.combatants = append(m.combatants, combatant{name: name, init: 0, hp: 10, maxHP: 10, conditions: []string{}})
					m.inputMode = InputModeAddInit
					m.textInput.Placeholder = "Enter initiative (or 'roll' for d20)"
					m.textInput.SetValue("")
				}
			} else if m.inputMode == InputModeAddInit {
				input := m.textInput.Value()
				if input == "roll" {
					m.combatants[len(m.combatants)-1].init = rand.Intn(20) + 1
				} else if init, err := strconv.Atoi(input); err == nil {
					m.combatants[len(m.combatants)-1].init = init
				}
				m.inputMode = InputModeAddName
				m.textInput.Placeholder = "Enter name (or 'done' to finish)"
				m.textInput.SetValue("")
				m.updateList()
			}
		case tea.KeyEsc:
			return m, func() tea.Msg { return switchModeMsg{"main"} }
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-ListHeightPadding)
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m initiativeTracker) View() string {
	var content strings.Builder
	content.WriteString("Initiative Tracker\n\n")
	content.WriteString("Combatants:\n")
	for _, c := range m.combatants {
		content.WriteString(fmt.Sprintf("%s (Init: %d, HP: %d/%d, Conditions: %s)\n", c.name, c.init, c.hp, c.maxHP, strings.Join(c.conditions, ", ")))
	}
	content.WriteString("\n")
	content.WriteString(m.textInput.View())
	content.WriteString("\n\nPress Esc to exit.")
	return content.String()
}

func (m *initiativeTracker) updateList() {
	// Sort combatants by initiative descending
	sort.Slice(m.combatants, func(i, j int) bool {
		return m.combatants[i].init > m.combatants[j].init
	})
	// Update list items if needed, but for now, just sort
}
