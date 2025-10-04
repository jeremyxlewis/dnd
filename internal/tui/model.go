package tui

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"dnd-cli/internal/data"
	"dnd-cli/internal/dice"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// mainModel represents the main command-line interface model.
// mainModel represents the main command-line interface model.
type mainModel struct {
	textInput    textinput.Model
	viewport     viewport.Model
	lastContent  string
	lastStyle    lipgloss.Style
	width        int
	height       int
	history      []string
	historyIndex int
	prevValue    string
	status       string
}

// topModel is the top-level model that manages switching between different sub-models.
type topModel struct {
	current tea.Model
	width   int
	height  int
}

// setWrappedContent sets the viewport content with word wrapping and optional styling.
func (m *mainModel) setWrappedContent(content string, style ...lipgloss.Style) {
	m.lastContent = content
	var s lipgloss.Style
	if len(style) > 0 {
		s = style[0]
		m.lastStyle = s
	} else {
		s = m.lastStyle
		if s.String() == "" { // if no last style, use default
			s = outputStyle
			m.lastStyle = s
		}
	}
	// For styled content, wrap to inner width first
	innerWidth := m.viewport.Width - 4 // border 2, padding 2
	wrappedContent := lipgloss.NewStyle().Width(innerWidth).Render(content)
	wrapped := s.Render(wrappedContent)
	m.viewport.SetContent(wrapped)
}

// newMainModel creates a new instance of the main TUI model with the given dimensions.
func newMainModel(width, height int) mainModel {
	ti := textinput.New()
	ti.Placeholder = "Type something..."
	ti.Focus()
	ti.CharLimit = TextInputCharLimit
	tiWidth := width - 20
	if tiWidth < 20 {
		tiWidth = 20
	}
	if tiWidth > TextInputWidth {
		tiWidth = TextInputWidth
	}
	ti.Width = tiWidth
	ti.PromptStyle = focusedStyle
	ti.TextStyle = focusedStyle
	ti.Cursor.Style = cursorStyle

	vpWidth := width - ViewportWidthPadding
	if vpWidth < 1 {
		vpWidth = 1
	}
	vpHeight := height - ViewportHeightPadding
	if vpHeight < 1 {
		vpHeight = 1
	}
	vp := viewport.New(vpWidth, vpHeight) // account for border/padding

	return mainModel{
		textInput:    ti,
		viewport:     vp,
		width:        width,
		height:       height,
		history:      []string{},
		historyIndex: 0,
		prevValue:    "",
		status:       "Ready. Type 'help' or '?' for commands.",
	}
}

// NewModel creates the top-level TUI model with initial dimensions.
func NewModel() topModel {
	return topModel{current: newMainModel(DefaultWidth, DefaultHeight), width: DefaultWidth, height: DefaultHeight}
}

// getHelpText returns a formatted help text for the TUI.
func getHelpText() string {
	return `Available Commands:

Core Commands:
   roll <notation>     - Roll dice (e.g., roll 1d20, roll 2d6+3)

 Lookup Commands:
     search [query]      - Global fuzzy search across all categories (spells, monsters, items, races, backgrounds, classes, rules)
     spell [name]        - Browse/filter spell list or look up specific spell
     monster [name]      - Browse/filter monster list or look up specific monster
     item [name]         - Browse/filter item list or look up specific item
     race [name]         - Browse/filter race list or look up specific race
     background [name]   - Browse/filter background list or look up specific background
     class [name]        - Browse/filter class list or look up specific class
     rules [topic]       - Look up PHB rules (combat, conditions, ability checks, etc.)

 Character Management (Full PHB Support):
    char create         - Create a new character interactively in TUI (ability scores, all races/classes/backgrounds)
    char view <name>    - View a character's full details
    char levelup <name> - Level up a character with proper mechanics
    char hp <name> <action> <amount> - Manage HP (damage/heal/set)
    char spells <name> <action> <level> <amount> - Manage spell slots (use/restore)
    char inventory <name> <action> <item> - Manage inventory (add/remove)
    char condition <name> <action> <condition> - Manage conditions (add/remove)
    char edit <name> <field> <value> - Edit character details (alignment/backstory)

 Combat & DM Tools:
    combat              - Launch initiative tracker for managing combat turns, HP, and conditions

NPC Generation:
   npc [generate]      - Generate a random NPC

Other:
   help or ?           - Show this help message

  Keyboard Shortcuts:
     Ctrl+H              - Show help
     Ctrl+C              - Quit
     Esc                 - Exit or go back

  Commands:
     quit or exit        - Quit the TUI

  In lists, type to filter, use arrows to navigate, Enter to select, Esc to cancel.
 Use ↑/↓ to scroll output.`
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
		if msg.Type == -1 && len(msg.Runes) == 1 && msg.Runes[0] == '/' {
			m.textInput.SetValue("")
			return m, func() tea.Msg { return switchModeMsg{"fuzzy_global"} }
		}
		switch msg.Type {
		case tea.KeyEnter:
			input := m.textInput.Value()
			if input != "" {
				// Add to history if not duplicate of last
				if len(m.history) == 0 || m.history[len(m.history)-1] != input {
					m.history = append(m.history, input)
					// Limit to HistoryLimit
					if len(m.history) > HistoryLimit {
						m.history = m.history[1:]
					}
				}
				m.historyIndex = len(m.history)
			}
			if input == "help" || input == "?" {
				m.setWrappedContent(getHelpText())
				m.status = "Help displayed. Use ↑/↓ to scroll."
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
							m.setWrappedContent(fmt.Sprintf("Error: %v", err), errorStyle)
						} else {
							total, rolls := dr.Roll()
							content := fmt.Sprintf("Rolling %s\n\nRolls: %v\n\nTotal: %d", dr.Notation, rolls, total)
							m.setWrappedContent(content, rollStyle)
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
							m.setWrappedContent(getRandomSpellErrorMessage(name), errorStyle)
						} else {
							content := fmt.Sprintf("--- %s ---\n\n", spell.Name)
							order := []string{"Level", "School", "Casting Time", "Range", "Components", "Duration"}
							content += renderPropertiesTable(spell.Properties, order)
							content += fmt.Sprintf("\n\nDescription:\n%s\n\n", formatDescription(spell.Description))
							content += fmt.Sprintf("Source: %s (%s)\n", spell.Book, spell.Publisher)
							m.setWrappedContent(content, infoCardStyle)
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
							m.setWrappedContent(getRandomMonsterErrorMessage(name), errorStyle)
						} else {
							content := fmt.Sprintf("--- %s ---\n\n", monster.Name)
							content += fmt.Sprintf("Description:\n%s\n", formatDescription(monster.Description))
							m.setWrappedContent(content, infoCardStyle)
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
							m.setWrappedContent(getRandomItemErrorMessage(name), errorStyle)
						} else {
							content := fmt.Sprintf("--- %s ---\n\n", it.Name)
							content += fmt.Sprintf("Description:\n%s\n", formatDescription(it.Description))
							m.setWrappedContent(content, infoCardStyle)
						}
					}
				case "race":
					if len(args) < 2 {
						m.textInput.SetValue("")
						return m, func() tea.Msg { return switchModeMsg{"fuzzy_race"} }
					} else {
						name := strings.Join(args[1:], " ")
						species, err := data.GetSpeciesByName(name)
						if err != nil {
							m.setWrappedContent(getRandomSpeciesErrorMessage(name), errorStyle)
						} else {
							content := fmt.Sprintf("--- %s ---\n\n", species.Name)
							content += fmt.Sprintf("Description:\n%s\n", formatDescription(species.Description))
							m.setWrappedContent(content, infoCardStyle)
						}
					}
				case "background":
					if len(args) < 2 {
						m.textInput.SetValue("")
						return m, func() tea.Msg { return switchModeMsg{"fuzzy_background"} }
					} else {
						name := strings.Join(args[1:], " ")
						background, err := data.GetBackgroundByName(name)
						if err != nil {
							m.setWrappedContent(getRandomBackgroundErrorMessage(name), errorStyle)
						} else {
							content := fmt.Sprintf("--- %s ---\n\n", background.Name)
							content += fmt.Sprintf("Description:\n%s\n", formatDescription(background.Description))
							m.setWrappedContent(content, infoCardStyle)
						}
					}
				case "class":
					if len(args) < 2 {
						m.textInput.SetValue("")
						return m, func() tea.Msg { return switchModeMsg{"fuzzy_class"} }
					} else {
						name := strings.Join(args[1:], " ")
						class, err := data.GetClassByName(name)
						if err != nil {
							m.setWrappedContent(getRandomClassErrorMessage(name), errorStyle)
						} else {
							content := fmt.Sprintf("--- %s ---\n\n", class.Name)
							content += fmt.Sprintf("Description:\n%s\n", formatDescription(class.Description))
							m.setWrappedContent(content, infoCardStyle)
						}
					}
				case "rules":
					if len(args) < 2 {
						m.textInput.SetValue("")
						return m, func() tea.Msg { return switchModeMsg{"fuzzy_rules"} }
					} else {
						topic := strings.ToLower(strings.Join(args[1:], " "))
						if desc, ok := rules[topic]; ok {
							content := fmt.Sprintf("--- %s ---\n\n%s", strings.Title(topic), formatDescription(desc))
							m.setWrappedContent(content, infoCardStyle)
						} else {
							m.setWrappedContent(fmt.Sprintf("Hark! No rules for '%s'. Available: combat, conditions, ability checks, initiative, actions.", topic), errorStyle)
						}
					}
				case "search":
					m.textInput.SetValue("")
					return m, func() tea.Msg { return switchModeMsg{"fuzzy_global"} }
				case "combat":
					m.textInput.SetValue("")
					return m, func() tea.Msg { return switchModeMsg{"initiative_tracker"} }
				case "npc":
					if len(args) < 2 || args[1] == "generate" {
						npc := data.GenerateNPC()
						content := fmt.Sprintf("--- Generated NPC ---\n\nName: %s\nSpecies: %s\nBackground: %s\n\nPersonality Trait: %s\n\nIdeal: %s\n\nBond: %s\n\nFlaw: %s\n\nBackstory: %s\n", npc.Name, npc.Species, npc.Background, npc.PersonalityTrait, npc.Ideal, npc.Bond, npc.Flaw, npc.Backstory)
						m.setWrappedContent(infoCardStyle.Render(content))
					} else {
						m.setWrappedContent("Unknown npc subcommand.", errorStyle)
					}
				case "char":
					if len(args) >= 2 && args[1] == "create" {
						m.textInput.SetValue("")
						return m, func() tea.Msg { return switchModeMsg{"char_create"} }
					} else {
						m.setWrappedContent("Usage: char create", errorStyle)
					}
				case "quit", "exit":
					return m, tea.Quit
				default:
					m.setWrappedContent(getRandomErrorMessage(), errorStyle)
				}
				m.textInput.SetValue("")
			}
		case tea.KeyUp:
			if len(m.history) > 0 && m.historyIndex > 0 {
				m.historyIndex--
				m.textInput.SetValue(m.history[m.historyIndex])
			}
		case tea.KeyDown:
			if m.historyIndex < len(m.history)-1 {
				m.historyIndex++
				m.textInput.SetValue(m.history[m.historyIndex])
			} else if m.historyIndex == len(m.history)-1 || m.historyIndex == len(m.history) {
				m.historyIndex = len(m.history)
				m.textInput.SetValue("")
			}
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyPgUp:
			m.viewport.HalfViewUp()
		case tea.KeyPgDown:
			m.viewport.HalfViewDown()
		}

	// We handle errors just like any other message
	case errMsg:
		// Handle error if needed
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	if m.textInput.Value() != m.prevValue {
		m.historyIndex = len(m.history)
	}
	m.prevValue = m.textInput.Value()
	return m, cmd
}

// View renders the UI.
func (m mainModel) View() string {
	viewportContent := m.viewport.View()

	promptSection := lipgloss.JoinVertical(lipgloss.Left,
		promptStyle.Render("What is thy command, adventurer?"),
		m.textInput.View(),
		quitStyle.Render("Press Ctrl+C to quit. Use ↑/↓ to scroll output."),
	)

	return lipgloss.JoinVertical(lipgloss.Left, viewportContent, promptSection)
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
			mm := newMainModel(m.width, m.height)
			m.current = mm
		case "char_create":
			m.current = newCharCreateModel(m.width, m.height)
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
		case "fuzzy_race":
			fm := newFuzzyModel("race")
			fm.list.SetSize(m.width, m.height-2)
			m.current = fm
		case "fuzzy_background":
			fm := newFuzzyModel("background")
			fm.list.SetSize(m.width, m.height-2)
			m.current = fm
		case "fuzzy_class":
			fm := newFuzzyModel("class")
			fm.list.SetSize(m.width, m.height-2)
			m.current = fm
		case "fuzzy_rules":
			fm := newFuzzyModel("rules")
			fm.list.SetSize(m.width, m.height-2)
			m.current = fm
		case "fuzzy_global":
			fm := newFuzzyModel("global")
			fm.list.SetSize(m.width, m.height-2)
			m.current = fm
		case "initiative_tracker":
			it := newInitiativeTracker(m.width, m.height)
			m.current = it
		}
		return m, nil
	case selectedMsg:
		mm := newMainModel(m.width, m.height)
		if msg.mode == "global" {
			// Parse "Category: Name"
			parts := strings.SplitN(msg.name, ": ", 2)
			if len(parts) != 2 {
				mm.setWrappedContent("Invalid global search result.", errorStyle)
			} else {
				category := strings.ToLower(parts[0])
				name := parts[1]
				displayItem(&mm, category, name)
			}
		} else {
			displayItem(&mm, msg.mode, msg.name)
		}
		m.current = mm
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlH:
			// Display help in main model
			if mm, ok := m.current.(*mainModel); ok {
				mm.setWrappedContent(getHelpText())
				mm.status = "Help displayed. Press Esc to return."
			}
			return m, nil
		default:
			// Pass to current model
			m.current, cmd = m.current.Update(msg)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update current model if it has size
		if mm, ok := m.current.(*mainModel); ok {
			mm.width = msg.Width
			mm.height = msg.Height
			mm.viewport.Width = msg.Width - ViewportWidthPadding
			mm.viewport.Height = msg.Height - ViewportHeightPadding
			tiWidth := msg.Width - 20
			if tiWidth < 20 {
				tiWidth = 20
			}
			if tiWidth > TextInputWidth {
				tiWidth = TextInputWidth
			}
			mm.textInput.Width = tiWidth
			// Re-wrap content if any
			if mm.lastContent != "" {
				mm.setWrappedContent(mm.lastContent, mm.lastStyle)
			}
		} else if fm, ok := m.current.(*fuzzyModel); ok {
			fm.list.SetSize(msg.Width, msg.Height-ListHeightPadding)
		} else if cm, ok := m.current.(*charCreateModel); ok {
			cm.width = msg.Width
			cm.height = msg.Height
			cm.list.SetSize(msg.Width, msg.Height-ListHeightPadding)
		} else if it, ok := m.current.(*initiativeTracker); ok {
			it.width = msg.Width
			it.height = msg.Height
			it.list.SetSize(msg.Width, msg.Height-ListHeightPadding)
		}
		return m, nil
	default:
		m.current, cmd = m.current.Update(msg)
		return m, cmd
	}
}

func (m topModel) View() string {
	view := m.current.View()
	// Fill to screen height with newlines to ensure UI is at top
	lineCount := strings.Count(view, "\n") + 1
	if lineCount < m.height {
		view += strings.Repeat("\n", m.height-lineCount)
	}
	return view
}

// switchModeMsg is used to switch between different TUI modes.
type switchModeMsg struct {
	mode string // The mode to switch to (e.g., "main", "fuzzy_spell")
}

// selectedMsg is sent when an item is selected from fuzzy finder.
type selectedMsg struct {
	mode string // The mode (e.g., "spell", "monster") or "global"
	name string // The selected item name
}

// displayItem displays the content for a given category and name in the main model.
// It handles fetching data and formatting for spells, monsters, items, races, backgrounds, classes, and rules.
func displayItem(mm *mainModel, category, name string) {
	switch category {
	case "spell":
		spell, err := data.GetSpellByName(name)
		if err != nil {
			mm.setWrappedContent(getRandomSpellErrorMessage(name), errorStyle)
		} else {
			content := fmt.Sprintf("--- %s ---\n\n", spell.Name)
			order := []string{"Level", "School", "Casting Time", "Range", "Components", "Duration"}
			content += renderPropertiesTable(spell.Properties, order)
			content += fmt.Sprintf("\n\nDescription:\n%s\n\n", formatDescription(spell.Description))
			content += fmt.Sprintf("Source: %s (%s)\n", spell.Book, spell.Publisher)
			mm.setWrappedContent(content, infoCardStyle)
		}
	case "monster":
		monster, err := data.GetMonsterByName(name)
		if err != nil {
			mm.setWrappedContent(getRandomMonsterErrorMessage(name), errorStyle)
		} else {
			content := fmt.Sprintf("--- %s ---\n\n", monster.Name)
			content += fmt.Sprintf("Description:\n%s\n", formatDescription(monster.Description))
			mm.setWrappedContent(content, infoCardStyle)
		}
	case "item":
		it, err := data.GetItemByName(name)
		if err != nil {
			mm.setWrappedContent(getRandomItemErrorMessage(name), errorStyle)
		} else {
			content := fmt.Sprintf("--- %s ---\n\n", it.Name)
			content += fmt.Sprintf("Description:\n%s\n", formatDescription(it.Description))
			mm.setWrappedContent(content, infoCardStyle)
		}
	case "race":
		species, err := data.GetSpeciesByName(name)
		if err != nil {
			mm.setWrappedContent(getRandomSpeciesErrorMessage(name), errorStyle)
		} else {
			content := fmt.Sprintf("--- %s ---\n\n", species.Name)
			content += fmt.Sprintf("Description:\n%s\n", formatDescription(species.Description))
			mm.setWrappedContent(content, infoCardStyle)
		}
	case "background":
		background, err := data.GetBackgroundByName(name)
		if err != nil {
			mm.setWrappedContent(getRandomBackgroundErrorMessage(name), errorStyle)
		} else {
			content := fmt.Sprintf("--- %s ---\n\n", background.Name)
			content += fmt.Sprintf("Description:\n%s\n", formatDescription(background.Description))
			mm.setWrappedContent(content, infoCardStyle)
		}
	case "class":
		class, err := data.GetClassByName(name)
		if err != nil {
			mm.setWrappedContent(getRandomClassErrorMessage(name), errorStyle)
		} else {
			content := fmt.Sprintf("--- %s ---\n\n", class.Name)
			content += fmt.Sprintf("Description:\n%s\n", formatDescription(class.Description))
			mm.setWrappedContent(content, infoCardStyle)
		}
	case "rules":
		if desc, ok := rules[name]; ok {
			content := fmt.Sprintf("--- %s ---\n\n%s", strings.Title(name), formatDescription(desc))
			mm.setWrappedContent(content, infoCardStyle)
		} else {
			mm.setWrappedContent("Unknown rules topic.", errorStyle)
		}
	default:
		mm.setWrappedContent("Unknown category.", errorStyle)
	}
}

// errMsg is a custom error type for our TUI (currently unused).
type errMsg error

// StartTUI runs the Bubble Tea application for the D&D CLI TUI.
func StartTUI() {
	config := LoadConfig()
	ApplyTheme(config.Theme)
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
