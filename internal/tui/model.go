package tui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"dnd-cli/internal/data"
	"dnd-cli/internal/dice"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()

	noStyle = lipgloss.NewStyle()

	headerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Background(lipgloss.Color("236")).Padding(0, 1).Bold(true)
	promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	quitStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	outputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Padding(0, 1)
)

type mainModel struct {
	textInput   textinput.Model
	viewport    viewport.Model
	lastContent string
}

type topModel struct {
	current tea.Model
	width   int
	height  int
}

// setWrappedContent sets the viewport content with word wrapping.
func (m *mainModel) setWrappedContent(content string) {
	m.lastContent = content
	wrapped := lipgloss.NewStyle().Width(m.viewport.Width - 2).Render(content)
	m.viewport.SetContent(wrapped)
}

// newMainModel creates a new instance of the main TUI model.
func newMainModel() mainModel {
	ti := textinput.New()
	ti.Placeholder = "Type something..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40
	ti.PromptStyle = focusedStyle
	ti.TextStyle = focusedStyle
	ti.Cursor.Style = cursorStyle

	vp := viewport.New(78, 20) // Initial size, adjusted on resize

	return mainModel{
		textInput: ti,
		viewport:  vp,
	}
}

// NewModel creates the top-level TUI model.
func NewModel() topModel {
	return topModel{current: newMainModel(), width: 80, height: 24}
}

// getHelpText returns a formatted help text for the TUI.
func getHelpText() string {
	return `Available Commands:

Core Commands:
  roll <notation>     - Roll dice (e.g., roll 1d20, roll 2d6+3)

Lookup Commands:
  spell [name]        - Browse/filter spell list or look up specific spell
  monster [name]      - Browse/filter monster list or look up specific monster
  item [name]         - Browse/filter item list or look up specific item

NPC Generation:
  npc [generate]      - Generate a random NPC

Other:
  help or ?           - Show this help message

In lists, type to filter, use arrows to navigate, Enter to select, Esc to cancel.
Press Esc or Ctrl+C to quit the TUI.`
}

// Init initializes the main TUI model. It can return a command to perform initial actions.
func (m mainModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the main model's state.
func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			input := m.textInput.Value()
			if input == "help" || input == "?" {
				m.setWrappedContent(getHelpText())
				m.textInput.SetValue("")
			} else if input != "" {
				args := strings.Fields(input)
				if len(args) == 0 {
					m.textInput.SetValue("")
					break
				}
				cmd := args[0]
				switch cmd {
				case "roll":
					if len(args) < 2 {
						m.setWrappedContent("Usage: roll <notation> (e.g., roll 1d20, roll 2d6+3)")
					} else {
						notation := args[1]
						dr, err := dice.ParseDiceNotation(notation)
						if err != nil {
							m.setWrappedContent(fmt.Sprintf("Error: %v", err))
						} else {
							total, rolls := dr.Roll()
							m.setWrappedContent(fmt.Sprintf("Rolling %s: %v -> Total: %d", dr.Notation, rolls, total))
						}
					}
				case "spell":
					if len(args) < 2 {
						m.textInput.SetValue("")
						return m, func() tea.Msg { return switchModeMsg{"fuzzy_spell"} }
					} else {
						name := strings.Join(args[1:], " ")
						spell, err := data.GetSpellByName(name)
						if err != nil {
							m.setWrappedContent(fmt.Sprintf("Hark! The spell '%s' is not etched in my scrolls.", name))
						} else {
							content := fmt.Sprintf("--- %s ---\n\n", spell.Name)
							order := []string{"Level", "School", "Casting Time", "Range", "Components", "Duration"}
							for _, key := range order {
								if v, ok := spell.Properties[key]; ok {
									content += fmt.Sprintf("%s: %v\n", key, v)
								}
							}
							content += fmt.Sprintf("\nDescription:\n%s\n\n", spell.Description)
							content += fmt.Sprintf("Source: %s (%s)\n", spell.Book, spell.Publisher)
							m.setWrappedContent(content)
						}
					}
				case "monster":
					if len(args) < 2 {
						m.textInput.SetValue("")
						return m, func() tea.Msg { return switchModeMsg{"fuzzy_monster"} }
					} else {
						name := strings.Join(args[1:], " ")
						monster, err := data.GetMonsterByName(name)
						if err != nil {
							m.setWrappedContent(fmt.Sprintf("Hark! The beast '%s' is not known in these lands.", name))
						} else {
							content := fmt.Sprintf("--- %s ---\n\n", monster.Name)
							content += fmt.Sprintf("Description:\n%s\n", strings.ReplaceAll(monster.Description, ". ", ".\n\n"))
							m.setWrappedContent(content)
						}
					}
				case "item":
					if len(args) < 2 {
						m.textInput.SetValue("")
						return m, func() tea.Msg { return switchModeMsg{"fuzzy_item"} }
					} else {
						name := strings.Join(args[1:], " ")
						it, err := data.GetItemByName(name)
						if err != nil {
							m.setWrappedContent(fmt.Sprintf("Hark! The item '%s' is not in my treasure hoard.", name))
						} else {
							content := fmt.Sprintf("--- %s ---\n\n", it.Name)
							content += fmt.Sprintf("Description:\n%s\n", strings.ReplaceAll(it.Description, ". ", ".\n\n"))
							m.setWrappedContent(content)
						}
					}
				case "npc":
					if len(args) < 2 || args[1] == "generate" {
						npc := data.GenerateNPC()
						content := fmt.Sprintf("--- Generated NPC ---\n\nName: %s\nSpecies: %s\nBackground: %s\n\nPersonality Trait: %s\n\nIdeal: %s\n\nBond: %s\n\nFlaw: %s\n\nBackstory: %s\n", npc.Name, npc.Species, npc.Background, npc.PersonalityTrait, npc.Ideal, npc.Bond, npc.Flaw, npc.Backstory)
						m.setWrappedContent(content)
					} else {
						m.setWrappedContent("Unknown npc subcommand.")
					}
				default:
					m.setWrappedContent("Hark! That command eludes my arcane senses. Type 'help' for available incantations.")
				}
				m.textInput.SetValue("")
			}
		case tea.KeyUp:
			m.viewport.LineUp(1)
		case tea.KeyDown:
			m.viewport.LineDown(1)
		case tea.KeyPgUp:
			m.viewport.HalfViewUp()
		case tea.KeyPgDown:
			m.viewport.HalfViewDown()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width - 2
		m.viewport.Height = msg.Height - 3
		if m.lastContent != "" {
			m.setWrappedContent(m.lastContent)
		}

	// We handle errors just like any other message
	case errMsg:
		// Handle error if needed
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the UI.
func (m mainModel) View() string {
	s := outputStyle.Render(m.viewport.View()) + "\n"

	s += promptStyle.Render("What is thy command, adventurer?") + "\n"
	s += m.textInput.View()
	s += quitStyle.Render("\nPress Esc or Ctrl+C to quit. Use ↑/↓ to scroll output.")

	return s
}

// topModel methods

func (m topModel) Init() tea.Cmd {
	return m.current.Init()
}

func (m topModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case switchModeMsg:
		switch msg.mode {
		case "main":
			mm := newMainModel()
			mm.viewport.Width = m.width - 2
			mm.viewport.Height = m.height - 3
			m.current = mm
		case "char_create":
			m.current = newCharCreateModel()
		case "fuzzy_spell":
			fm := newFuzzyModel("spell")
			fm.list.SetSize(m.width, m.height-2)
			m.current = fm
		case "fuzzy_monster":
			fm := newFuzzyModel("monster")
			fm.list.SetSize(m.width, m.height-2)
			m.current = fm
		case "fuzzy_item":
			fm := newFuzzyModel("item")
			fm.list.SetSize(m.width, m.height-2)
			m.current = fm
		}
		return m, nil
	case selectedMsg:
		mm := newMainModel()
		mm.viewport.Width = m.width - 2
		mm.viewport.Height = m.height - 3
		switch msg.mode {
		case "spell":
			spell, err := data.GetSpellByName(msg.name)
			if err != nil {
				mm.setWrappedContent(fmt.Sprintf("Hark! The spell '%s' is not etched in my scrolls.", msg.name))
			} else {
				content := fmt.Sprintf("--- %s ---\n\n", spell.Name)
				order := []string{"Level", "School", "Casting Time", "Range", "Components", "Duration"}
				for _, key := range order {
					if v, ok := spell.Properties[key]; ok {
						content += fmt.Sprintf("%s: %v\n", key, v)
					}
				}
				content += fmt.Sprintf("\nDescription:\n%s\n\n", spell.Description)
				content += fmt.Sprintf("Source: %s (%s)\n", spell.Book, spell.Publisher)
				mm.setWrappedContent(content)
			}
		case "monster":
			monster, err := data.GetMonsterByName(msg.name)
			if err != nil {
				mm.setWrappedContent(fmt.Sprintf("Hark! The beast '%s' is not known in these lands.", msg.name))
			} else {
				content := fmt.Sprintf("--- %s ---\n\n", monster.Name)
				content += fmt.Sprintf("Description:\n%s\n", strings.ReplaceAll(monster.Description, ". ", ".\n\n"))
				mm.setWrappedContent(content)
			}
		case "item":
			it, err := data.GetItemByName(msg.name)
			if err != nil {
				mm.setWrappedContent(fmt.Sprintf("Hark! The item '%s' is not in my treasure hoard.", msg.name))
			} else {
				content := fmt.Sprintf("--- %s ---\n\n", it.Name)
				content += fmt.Sprintf("Description:\n%s\n", strings.ReplaceAll(it.Description, ". ", ".\n\n"))
				mm.setWrappedContent(content)
			}
		}
		m.current = mm
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update current model if it has size
		if mm, ok := m.current.(*mainModel); ok {
			mm.viewport.Width = msg.Width - 2
			mm.viewport.Height = msg.Height - 3
		} else if fm, ok := m.current.(*fuzzyModel); ok {
			fm.list.SetSize(msg.Width, msg.Height-2)
		}
		return m, nil
	default:
		m.current, cmd = m.current.Update(msg)
		return m, cmd
	}
}

func (m topModel) View() string {
	return m.current.View()
}

// switchModeMsg is used to switch between different TUI modes.
type switchModeMsg struct {
	mode string
}

// selectedMsg is sent when an item is selected from fuzzy finder.
type selectedMsg struct {
	mode string
	name string
}

// listItem represents a list item for fuzzy finder.
type listItem struct {
	title string
}

func (i listItem) FilterValue() string { return i.title }

// customDelegate for rendering list items.
type customDelegate struct{}

func (d customDelegate) Height() int { return 1 }

func (d customDelegate) Spacing() int { return 0 }

func (d customDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d customDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	li, ok := item.(listItem)
	if !ok {
		return
	}
	if index == m.Index() {
		fmt.Fprint(w, "> "+li.title)
	} else {
		fmt.Fprint(w, "  "+li.title)
	}
}

// charCreateModel handles the character creation mode.
type charCreateModel struct{}

func newCharCreateModel() charCreateModel {
	return charCreateModel{}
}

func (m charCreateModel) Init() tea.Cmd {
	return nil
}

func (m charCreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			return m, func() tea.Msg { return switchModeMsg{"main"} }
		}
	}
	return m, nil
}

func (m charCreateModel) View() string {
	return "Character Creation Mode\n\nThis is a placeholder. Press Esc to return to main."
}

// fuzzyModel handles fuzzy finding for spells, monsters, items.
type fuzzyModel struct {
	list  list.Model
	mode  string
	width int
}

func newFuzzyModel(mode string) fuzzyModel {
	var items []list.Item
	seen := make(map[string]bool)
	switch mode {
	case "spell":
		for _, s := range data.AllSpells {
			if !seen[s.Name] {
				items = append(items, listItem{title: s.Name})
				seen[s.Name] = true
			}
		}
	case "monster":
		for _, m := range data.AllMonsters {
			if !seen[m.Name] {
				items = append(items, listItem{title: m.Name})
				seen[m.Name] = true
			}
		}
	case "item":
		for _, i := range data.AllItems {
			if !seen[i.Name] {
				items = append(items, listItem{title: i.Name})
				seen[i.Name] = true
			}
		}
	}
	l := list.New(items, customDelegate{}, 80, 20) // initial size
	l.Title = "Select " + mode
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)
	return fuzzyModel{list: l, mode: mode, width: 80}
}

func (m fuzzyModel) Init() tea.Cmd {
	return nil
}

func (m fuzzyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		height := msg.Height - 4
		if height < 10 {
			height = 10
		}
		m.list.SetSize(msg.Width, height)
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			selected := m.list.SelectedItem()
			if selected != nil {
				name := selected.(listItem).title
				return m, func() tea.Msg { return selectedMsg{mode: m.mode, name: name} }
			}
		} else if msg.Type == tea.KeyEsc {
			return m, func() tea.Msg { return switchModeMsg{"main"} }
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m fuzzyModel) View() string {
	return m.list.View()
}

// errMsg is a custom error type for our TUI.
type errMsg error

// StartTUI runs the Bubble Tea application.
func StartTUI() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
