package tui

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"

	"dnd-cli/internal/character"
	"dnd-cli/internal/data"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// charCreateModel handles the character creation mode.
type charCreateModel struct {
	step          int // 0: name, 1: alignment, 2: player, 3: level, 4: score method, 5: scores, 6: species, 7: species info, 8: class, 9: class info, 10: background, 11: background info, 12: proficiencies, 13: equipment, 14: spellcasting, 15: confirm
	name          string
	alignment     string
	player        string
	level         int
	scores        [6]int // str, dex, con, int, wis, cha
	scoreMethod   string
	species       string
	class         string
	background    string
	proficiencies []string
	equipment     []string
	spells        []string
	textInput     textinput.Model
	list          list.Model
	width         int
	height        int
}

func newCharCreateModel(width, height int) charCreateModel {
	ti := textinput.New()
	ti.Placeholder = "Enter character name"
	ti.Focus()
	ti.CharLimit = TextInputCharLimit

	return charCreateModel{
		step:          StepName,
		level:         1,
		textInput:     ti,
		width:         width,
		height:        height,
		proficiencies: []string{},
		equipment:     []string{},
		spells:        []string{},
	}
}

func (m charCreateModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m charCreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.step == StepName {
				m.name = m.textInput.Value()
				if m.name == "" {
					// Stay on name
					break
				}
				// Check if exists
				charFilePath, err := character.GetCharacterFilePath(m.name)
				if err == nil {
					if _, err := os.Stat(charFilePath); err == nil {
						// Exists, stay
						break
					}
				}
				// Proceed to alignment
				m.step = StepAlignment
				m.setupAlignmentList()
			} else if m.step == StepPlayer {
				m.player = m.textInput.Value()
				// Proceed to level
				m.step = StepLevel
				m.textInput.Placeholder = "Enter starting level (default 1)"
				m.textInput.SetValue("1")
			} else if m.step == StepLevel {
				levelStr := m.textInput.Value()
				if levelStr == "" {
					m.level = 1
				} else {
					if lvl, err := strconv.Atoi(levelStr); err == nil && lvl > 0 {
						m.level = lvl
					} else {
						m.level = 1
					}
				}
				// Proceed to score method
				m.step = StepScoreMethod
				m.setupScoreMethodList()
			} else if m.step == StepScores {
				// Accept scores
				m.step = StepSpecies
				m.setupSpeciesList()
			} else if m.step == StepAlignment || (m.step >= StepScoreMethod && m.step <= StepBackground && m.step != StepSpeciesInfo && m.step != StepClassInfo && m.step != StepBackgroundInfo) {
				selected := m.list.SelectedItem()
				if selected != nil {
					name := selected.(listItem).title
					switch m.step {
					case StepAlignment:
						m.alignment = name
						m.step = StepPlayer
						m.textInput.Placeholder = "Enter player/campaign details (optional)"
						m.textInput.SetValue("")
					case StepScoreMethod:
						m.scoreMethod = name
						m.step = StepScores
						m.generateScores()
					case StepSpecies:
						m.species = name
						m.step = StepSpeciesInfo
					case StepClass:
						m.class = name
						m.step = StepClassInfo
					case StepBackground:
						m.background = name
						m.step = StepBackgroundInfo
					}
				}
			} else if m.step == StepSpeciesInfo {
				// Apply racial ASIs temporarily
				m.applyRacialASIs()
				m.step = StepClass
				m.setupClassList()
			} else if m.step == StepClassInfo {
				// Display class info, but for now just proceed
				m.step = StepBackground
				m.setupBackgroundList()
			} else if m.step == StepBackgroundInfo {
				// Apply background proficiencies
				m.applyBackgroundProficiencies()
				m.step = StepProficiencies
			} else if m.step == StepProficiencies {
				m.step = StepEquipment
			} else if m.step == StepEquipment {
				m.step = StepConfirm
			} else if m.step == StepConfirm {
				// Confirm and create
				if m.name == "" {
					// Stay, perhaps show error
					break
				}
				if m.species == "" || m.class == "" || m.background == "" {
					break
				}
				// Validate scores
				valid := true
				for _, score := range m.scores {
					if score < 1 || score > 20 {
						valid = false
						break
					}
				}
				if !valid {
					break
				}
				// Create character
				char := character.NewCharacter(m.name, m.species, m.class, m.background, m.alignment, m.level, m.scores[0], m.scores[1], m.scores[2], m.scores[3], m.scores[4], m.scores[5])
				char.ApplyRacialTraits()
				char.ApplyClassTraits()
				char.ApplyBackgroundTraits()
				charFilePath, err := character.GetCharacterFilePath(m.name)
				if err == nil {
					character.SaveCharacter(char, charFilePath)
				}
				// Return to main
				return m, func() tea.Msg { return switchModeMsg{"main"} }
			}
		case tea.KeyEsc:
			if m.step == StepName {
				return m, func() tea.Msg { return switchModeMsg{"main"} }
			} else if m.step == StepScores {
				m.step = StepScoreMethod
				m.setupScoreMethodList()
			} else {
				m.step--
				if m.step == StepName {
					m.textInput.SetValue(m.name)
					m.textInput.Placeholder = "Enter character name"
					m.textInput.Focus()
				} else if m.step == StepAlignment {
					m.setupAlignmentList()
				} else if m.step == StepPlayer {
					m.textInput.SetValue(m.player)
					m.textInput.Placeholder = "Enter player/campaign details (optional)"
					m.textInput.Focus()
				} else if m.step == StepLevel {
					m.textInput.SetValue(strconv.Itoa(m.level))
					m.textInput.Placeholder = "Enter starting level (default 1)"
					m.textInput.Focus()
				} else if m.step == StepSpeciesInfo {
					m.step = StepSpecies
					m.setupSpeciesList()
				} else if m.step == StepClassInfo {
					m.step = StepClass
					m.setupClassList()
				} else if m.step == StepBackgroundInfo {
					m.step = StepBackground
					m.setupBackgroundList()
				} else if m.step == StepProficiencies {
					m.step = StepBackgroundInfo
				} else if m.step == StepEquipment {
					m.step = StepProficiencies
				} else {
					m.setupListForStep()
				}
			}
		}

	}

	if m.step == StepName || m.step == StepPlayer || m.step == StepLevel {
		m.textInput, cmd = m.textInput.Update(msg)
	} else if m.step == StepAlignment || (m.step >= StepScoreMethod && m.step <= StepBackground) {
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m *charCreateModel) setupScoreMethodList() {
	items := createListItems([]string{"Standard Array", "Roll 4d6 Drop Lowest", "Point Buy"})
	l := list.New(items, customDelegate{}, m.width, m.height-ListHeightPadding)
	l.KeyMap.Quit = key.NewBinding(key.WithDisabled())
	l.Title = "Select Ability Score Method"
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)
	l.Styles.Title = headerStyle
	l.Styles.FilterPrompt = focusedStyle
	l.Styles.FilterCursor = cursorStyle
	m.list = l
}

func (m *charCreateModel) setupAlignmentList() {
	items := createListItems([]string{"Lawful Good", "Neutral Good", "Chaotic Good", "Lawful Neutral", "Neutral", "Chaotic Neutral", "Lawful Evil", "Neutral Evil", "Chaotic Evil"})
	l := list.New(items, customDelegate{}, m.width, m.height-ListHeightPadding)
	l.KeyMap.Quit = key.NewBinding(key.WithDisabled())
	l.Title = "Select Alignment"
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)
	l.Styles.Title = headerStyle
	l.Styles.FilterPrompt = focusedStyle
	l.Styles.FilterCursor = cursorStyle
	m.list = l
}

func (m *charCreateModel) generateScores() {
	if m.scoreMethod == "Standard Array" || m.scoreMethod == "Point Buy" {
		m.scores = [6]int{15, 14, 13, 12, 10, 8}
	} else if m.scoreMethod == "Roll 4d6 Drop Lowest" {
		for i := range m.scores {
			rolls := make([]int, 4)
			for j := range rolls {
				rolls[j] = rand.Intn(6) + 1
			}
			// sort and drop lowest
			sort.Ints(rolls)
			sum := rolls[1] + rolls[2] + rolls[3]
			m.scores[i] = sum
		}
	}
	// Validate scores are reasonable (3-18)
	for i, score := range m.scores {
		if score < 3 || score > 18 {
			m.scores[i] = 10 // default
		}
	}
}

func (m *charCreateModel) setupSpeciesList() {
	titles := getUniqueTitles(data.AllSpecies, func(s data.Species) string { return s.Name })
	items := createListItems(titles)
	l := list.New(items, customDelegate{}, m.width, m.height-ListHeightPadding)
	l.KeyMap.Quit = key.NewBinding(key.WithDisabled())
	l.Title = "Select Species"
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)
	l.Styles.Title = headerStyle
	l.Styles.FilterPrompt = focusedStyle
	l.Styles.FilterCursor = cursorStyle
	m.list = l
}

func (m *charCreateModel) setupClassList() {
	titles := getUniqueTitles(data.AllClasses, func(c data.Class) string { return c.Name })
	items := createListItems(titles)
	l := list.New(items, customDelegate{}, m.width, m.height-ListHeightPadding)
	l.KeyMap.Quit = key.NewBinding(key.WithDisabled())
	l.Title = "Select Class"
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)
	l.Styles.Title = headerStyle
	l.Styles.FilterPrompt = focusedStyle
	l.Styles.FilterCursor = cursorStyle
	m.list = l
}

func (m *charCreateModel) setupBackgroundList() {
	titles := getUniqueTitles(data.AllBackgrounds, func(b data.Background) string { return b.Name })
	items := createListItems(titles)
	l := list.New(items, customDelegate{}, m.width, m.height-ListHeightPadding)
	l.KeyMap.Quit = key.NewBinding(key.WithDisabled())
	l.Title = "Select Background"
	l.SetFilteringEnabled(true)
	l.SetShowFilter(true)
	l.Styles.Title = headerStyle
	l.Styles.FilterPrompt = focusedStyle
	l.Styles.FilterCursor = cursorStyle
	m.list = l
}

func (m *charCreateModel) applyRacialASIs() {
	switch m.species {
	case "Human":
		m.scores[0]++ // str
		m.scores[1]++ // dex
		m.scores[2]++ // con
		m.scores[3]++ // int
		m.scores[4]++ // wis
		m.scores[5]++ // cha
	case "Elf", "High Elf", "Wood Elf", "Dark Elf":
		m.scores[1] += 2 // dex
		if m.species == "High Elf" {
			m.scores[3]++ // int
		} else if m.species == "Wood Elf" {
			m.scores[4]++ // wis
		} else if m.species == "Dark Elf" {
			m.scores[5]++ // cha
		}
	case "Dwarf", "Hill Dwarf", "Mountain Dwarf":
		m.scores[2] += 2 // con
		if m.species == "Hill Dwarf" {
			m.scores[4]++ // wis
		} else if m.species == "Mountain Dwarf" {
			m.scores[0]++ // str
		}
	case "Halfling", "Lightfoot", "Stout":
		m.scores[1] += 2 // dex
		if m.species == "Lightfoot" {
			m.scores[5]++ // cha
		} else if m.species == "Stout" {
			m.scores[2]++ // con
		}
	case "Half-Elf":
		m.scores[5] += 2 // cha
		// +1 to two others, placeholder +1 str +1 dex
		m.scores[0]++
		m.scores[1]++
	case "Half-Orc":
		m.scores[0] += 2 // str
		m.scores[2]++    // con
	case "Tiefling":
		m.scores[3]++    // int
		m.scores[5] += 2 // cha
	case "Dragonborn":
		m.scores[0] += 2 // str
		m.scores[5]++    // cha
	case "Gnome", "Forest Gnome", "Rock Gnome":
		m.scores[3] += 2 // int
		if m.species == "Forest Gnome" {
			m.scores[1]++ // dex
		} else if m.species == "Rock Gnome" {
			m.scores[2]++ // con
		}
	case "Tabaxi":
		m.scores[1] += 2 // dex
		m.scores[5]++    // cha
	case "Genasi", "Air Genasi", "Earth Genasi", "Fire Genasi", "Water Genasi":
		m.scores[2] += 2 // con
		if m.species == "Air Genasi" {
			m.scores[1]++ // dex
		} else if m.species == "Earth Genasi" {
			m.scores[0]++ // str
		} else if m.species == "Fire Genasi" {
			m.scores[3]++ // int
		} else if m.species == "Water Genasi" {
			m.scores[4]++ // wis
		}
	case "Aarakocra":
		m.scores[1] += 2 // dex
		m.scores[4]++    // wis
	case "Aasimar", "Protector Aasimar", "Scourge Aasimar", "Fallen Aasimar":
		m.scores[4]++    // wis
		m.scores[5] += 2 // cha
		if m.species == "Fallen Aasimar" {
			m.scores[0]++ // str
		}
	case "Bugbear":
		m.scores[0] += 2 // str
		m.scores[1]++    // dex
	case "Centaur":
		m.scores[0] += 2 // str
		m.scores[4]++    // wis
	case "Changeling":
		m.scores[5] += 2 // cha
	case "Duergar":
		m.scores[0] += 2 // str
		m.scores[2]++    // con
	case "Eladrin":
		m.scores[1] += 2 // dex
		m.scores[5]++    // cha
	case "Fairy":
		m.scores[1] += 2 // dex
		m.scores[5]++    // cha
	case "Firbolg":
		m.scores[4] += 2 // wis
		m.scores[0]++    // str
	case "Gith", "Githyanki", "Githzerai":
		m.scores[3] += 2 // int
		if m.species == "Githyanki" {
			m.scores[0]++ // str
		} else if m.species == "Githzerai" {
			m.scores[4]++ // wis
		}
	case "Goblin":
		m.scores[1] += 2 // dex
		m.scores[2]++    // con
	case "Goliath":
		m.scores[0] += 2 // str
		m.scores[2]++    // con
	case "Harengon":
		m.scores[1] += 2 // dex
		m.scores[4]++    // wis
	case "Hobgoblin":
		m.scores[2] += 2 // con
		m.scores[3]++    // int
	case "Kenku":
		m.scores[1] += 2 // dex
		m.scores[4]++    // wis
	case "Kobold":
		m.scores[1] += 2 // dex
		m.scores[0]--    // str
	case "Lizardfolk":
		m.scores[2] += 2 // con
		m.scores[4]++    // wis
	case "Minotaur":
		m.scores[0] += 2 // str
		m.scores[2]++    // con
	case "Orc":
		m.scores[0] += 2 // str
		m.scores[2]++    // con
	case "Satyr":
		m.scores[5] += 2 // cha
		m.scores[1]++    // dex
	case "Sea Elf":
		m.scores[2] += 2 // con
		m.scores[1]++    // dex
	case "Shifter", "Beasthide", "Cliffwalk", "Longstride", "Longtooth", "Razorclaw", "Wildhunt":
		m.scores[1] += 2 // dex
		m.scores[0]++    // str
		if m.species == "Beasthide" {
			m.scores[2]++ // con
		} else if m.species == "Wildhunt" {
			m.scores[4]++ // wis
		}
	case "Tortle":
		m.scores[0] += 2 // str
		m.scores[4]++    // wis
	case "Triton":
		m.scores[0]++ // str
		m.scores[2]++ // con
		m.scores[5]++ // cha
	case "Yuan-ti Pureblood":
		m.scores[5] += 2 // cha
		m.scores[3]++    // int
	}
}

func (m *charCreateModel) applyBackgroundProficiencies() {
	switch m.background {
	case "Acolyte":
		m.proficiencies = append(m.proficiencies, "Insight", "Religion")
	case "Charlatan":
		m.proficiencies = append(m.proficiencies, "Deception", "Sleight of Hand")
	case "Criminal":
		m.proficiencies = append(m.proficiencies, "Deception", "Stealth")
	case "Entertainer":
		m.proficiencies = append(m.proficiencies, "Acrobatics", "Performance")
	case "Folk Hero":
		m.proficiencies = append(m.proficiencies, "Animal Handling", "Survival")
	case "Guild Artisan":
		m.proficiencies = append(m.proficiencies, "Insight", "Persuasion")
	case "Hermit":
		m.proficiencies = append(m.proficiencies, "Medicine", "Religion")
	case "Noble":
		m.proficiencies = append(m.proficiencies, "History", "Persuasion")
	case "Outlander":
		m.proficiencies = append(m.proficiencies, "Athletics", "Survival")
	case "Sage":
		m.proficiencies = append(m.proficiencies, "Arcana", "History")
	case "Sailor":
		m.proficiencies = append(m.proficiencies, "Athletics", "Perception")
	case "Soldier":
		m.proficiencies = append(m.proficiencies, "Athletics", "Intimidation")
	case "Urchin":
		m.proficiencies = append(m.proficiencies, "Sleight of Hand", "Stealth")
	}
}

func getClassDescription(className string) string {
	for _, c := range data.AllClasses {
		if c.Name == className {
			return c.Description
		}
	}
	return "Description not found."
}

func getSpeciesDescription(speciesName string) string {
	for _, s := range data.AllSpecies {
		if s.Name == speciesName {
			return s.Description
		}
	}
	return "Description not found."
}

func getBackgroundDescription(backgroundName string) string {
	for _, b := range data.AllBackgrounds {
		if b.Name == backgroundName {
			return b.Description
		}
	}
	return "Description not found."
}

func (m *charCreateModel) setupListForStep() {
	switch m.step {
	case StepScoreMethod:
		m.setupScoreMethodList()
	case StepSpecies:
		m.setupSpeciesList()
	case StepClass:
		m.setupClassList()
	case StepBackground:
		m.setupBackgroundList()
	}
}

func (m charCreateModel) View() string {
	switch m.step {
	case StepName:
		return viewStyle.Render(fmt.Sprintf("Character Creation - Enter Name\n\n%s\n\nPress Enter to continue, Esc to cancel.", m.textInput.View()))
	case StepAlignment:
		return viewStyle.Render(fmt.Sprintf("Character Creation - Select Alignment\n\n%s\n\nType / to search, ↑↓ or jk to navigate, Enter to select, Esc to go back.", m.list.View()))
	case StepPlayer:
		return viewStyle.Render(fmt.Sprintf("Character Creation - Enter Player/Campaign Details\n\n%s\n\nPress Enter to continue, Esc to go back.", m.textInput.View()))
	case StepLevel:
		return viewStyle.Render(fmt.Sprintf("Character Creation - Enter Starting Level\n\n%s\n\nPress Enter to continue, Esc to go back.", m.textInput.View()))
	case StepScoreMethod:
		return viewStyle.Render(fmt.Sprintf("Character Creation - Select Score Method\n\n%s\n\nType / to search, ↑↓ or jk to navigate, Enter to select, Esc to go back.", m.list.View()))
	case StepScores:
		scoreNames := []string{"STR", "DEX", "CON", "INT", "WIS", "CHA"}
		scoreDisplay := ""
		for i, score := range m.scores {
			scoreDisplay += fmt.Sprintf("%s: %d\n", scoreNames[i], score)
		}
		return viewStyle.Render(fmt.Sprintf("Character Creation - Ability Scores (%s)\n\n%s\n\nPress Enter to accept, Esc to go back.", m.scoreMethod, scoreDisplay))
	case StepSpecies:
		return viewStyle.Render(fmt.Sprintf("Character Creation - Select Species\n\n%s\n\nType / to search, ↑↓ or jk to navigate, Enter to select, Esc to go back.", m.list.View()))
	case StepSpeciesInfo:
		desc := getSpeciesDescription(m.species)
		return viewStyle.Render(fmt.Sprintf("Character Creation - Species: %s\n\n%s\n\nRacial ASIs applied temporarily.\n\nPress Enter to continue, Esc to go back.", m.species, desc))
	case StepClass:
		return viewStyle.Render(fmt.Sprintf("Character Creation - Select Class\n\n%s\n\nType / to search, ↑↓ or jk to navigate, Enter to select, Esc to go back.", m.list.View()))
	case StepClassInfo:
		desc := getClassDescription(m.class)
		return viewStyle.Render(fmt.Sprintf("Character Creation - Class: %s\n\n%s\n\nPress Enter to continue, Esc to go back.", m.class, desc))
	case StepBackground:
		return viewStyle.Render(fmt.Sprintf("Character Creation - Select Background\n\n%s\n\nType / to search, ↑↓ or jk to navigate, Enter to select, Esc to go back.", m.list.View()))
	case StepBackgroundInfo:
		desc := getBackgroundDescription(m.background)
		return viewStyle.Render(fmt.Sprintf("Character Creation - Background: %s\n\n%s\n\nBackground proficiencies applied.\n\nPress Enter to continue, Esc to go back.", m.background, desc))
	case StepProficiencies:
		profStr := strings.Join(m.proficiencies, ", ")
		return viewStyle.Render(fmt.Sprintf("Character Creation - Proficiencies\n\nApplied Proficiencies: %s\n\nPress Enter to continue, Esc to go back.", profStr))
	case StepEquipment:
		equipStr := strings.Join(m.equipment, ", ")
		return viewStyle.Render(fmt.Sprintf("Character Creation - Equipment\n\nStarting Equipment: %s\n\nPress Enter to continue, Esc to go back.", equipStr))
	case StepConfirm:
		scoreNames := []string{"STR", "DEX", "CON", "INT", "WIS", "CHA"}
		scoreDisplay := ""
		for i, score := range m.scores {
			scoreDisplay += fmt.Sprintf("%s: %d ", scoreNames[i], score)
		}
		return viewStyle.Render(fmt.Sprintf("Character Creation - Confirm\n\nName: %s\nAlignment: %s\nPlayer: %s\nLevel: %d\nScores: %s\nSpecies: %s\nClass: %s\nBackground: %s\n\nPress Enter to create, Esc to go back.", m.name, m.alignment, m.player, m.level, scoreDisplay, m.species, m.class, m.background))
	default:
		return viewStyle.Render("Error")
	}
}
