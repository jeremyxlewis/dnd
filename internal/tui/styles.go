package tui

import "github.com/charmbracelet/lipgloss"

var isDark = lipgloss.HasDarkBackground()

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "11", Dark: "220"}) // Gold
	blurStyle    = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "240", Dark: "250"})
	cursorStyle  = focusedStyle

	noStyle = lipgloss.NewStyle()

	headerStyle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "11", Dark: "220"}).Background(lipgloss.AdaptiveColor{Light: "236", Dark: "236"}).Padding(0, 1).Bold(true) // Gold on gray
	promptStyle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "2", Dark: "2"}).Bold(true)                                                                                // Bright green
	quitStyle       = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "244", Dark: "244"}).Italic(true)                                                                          // Light gray
	outputStyle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "15", Dark: "15"}).Padding(0, 1)                                                                           // Bright white
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "1", Dark: "1"}).Padding(0, 1)                                                                             // Bright red
	viewportStyle   = lipgloss.NewStyle()
	selectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "11", Dark: "220"}).Bold(true)                                        // Gold
	unselectedStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "250"})                                                  // Light gray
	infoCardStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1).BorderForeground(lipgloss.AdaptiveColor{Light: "11", Dark: "220"}) // Gold border
	rollStyle       = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1).BorderForeground(lipgloss.AdaptiveColor{Light: "2", Dark: "2"})     // Green border
)

var errorMessages = []string{
	"Hark! That command eludes my arcane senses. Type 'help' for available incantations.",
	"Alas! Thy words are shrouded in mystery. Seek 'help' to unveil the secrets.",
	"By the gods! Such a command is unknown to me. Type 'help' for guidance.",
	"Fie! That incantation is not in my grimoire. Type 'help' to see the spells I know.",
	"Confusion reigns! Thy command is lost in the mists. 'Help' shall light the way.",
	"Oh no! My magical ears have failed me. Try 'help' for the proper enchantments.",
	"Zounds! That directive baffles even the wisest sages. 'Help' is thy ally.",
	"Egad! Such a command hath never graced my presence. 'Help' awaits thy call.",
	"Goodness me! Thy input is as enigmatic as a dragon's riddle. Seek 'help'.",
	"Heavens! That order is beyond my ken. 'Help' will reveal the path forward.",
}

var spellErrorMessages = []string{
	"Hark! The spell '%s' is not etched in my scrolls. Type 'help' for a list of spells.",
	"Alas! The incantation '%s' remains a mystery to me. Type 'help' for a list of spells.",
	"By the gods! Such a spell '%s' is unknown in these realms. Type 'help' for a list of spells.",
	"Fie! '%s' is not inscribed in my ancient tomes. Type 'help' for a list of spells.",
	"Confusion! The spell '%s' eludes my magical sight. Type 'help' for a list of spells.",
	"Oh no! '%s' is not among the spells I wield. Type 'help' for a list of spells.",
	"Zounds! '%s' baffles even the greatest wizards. Type 'help' for a list of spells.",
	"Egad! The spell '%s' hath never been cast here. Type 'help' for a list of spells.",
	"Goodness me! '%s' is as elusive as a shadow. Type 'help' for a list of spells.",
	"Heavens! '%s' is beyond my arcane knowledge. Type 'help' for a list of spells.",
}

var monsterErrorMessages = []string{
	"Hark! The beast '%s' is not known in these lands. Type 'help' for a list of monsters.",
	"Alas! The creature '%s' lurks not in my bestiaries. Type 'help' for a list of monsters.",
	"By the gods! Such a monster '%s' is unheard of. Type 'help' for a list of monsters.",
	"Fie! '%s' is not among the beasts I've encountered. Type 'help' for a list of monsters.",
	"Confusion! The monster '%s' hides from my gaze. Type 'help' for a list of monsters.",
	"Oh no! '%s' is not in the wilds I know. Type 'help' for a list of monsters.",
	"Zounds! '%s' confounds the bravest adventurers. Type 'help' for a list of monsters.",
	"Egad! The beast '%s' hath never crossed my path. Type 'help' for a list of monsters.",
	"Goodness me! '%s' is as mythical as a unicorn. Type 'help' for a list of monsters.",
	"Heavens! '%s' is unknown to mortal ken. Type 'help' for a list of monsters.",
}

var itemErrorMessages = []string{
	"Hark! The item '%s' is not in my treasure hoard. Type 'help' for a list of items.",
	"Alas! The artifact '%s' is lost to the ages. Type 'help' for a list of items.",
	"By the gods! Such an item '%s' is not in my vaults. Type 'help' for a list of items.",
	"Fie! '%s' is not among my glittering treasures. Type 'help' for a list of items.",
	"Confusion! The item '%s' evades my collection. Type 'help' for a list of items.",
	"Oh no! '%s' is not in my adventurer's pack. Type 'help' for a list of items.",
	"Zounds! '%s' baffles the greediest dragons. Type 'help' for a list of items.",
	"Egad! The item '%s' hath never been hoarded. Type 'help' for a list of items.",
	"Goodness me! '%s' is as rare as a philosopher's stone. Type 'help' for a list of items.",
	"Heavens! '%s' is beyond my material grasp. Type 'help' for a list of items.",
}

var speciesErrorMessages = []string{
	"Hark! The species '%s' is not known in these realms. Type 'help' for a list of species.",
	"Alas! The race '%s' lurks not in my tomes. Type 'help' for a list of species.",
	"By the gods! Such a species '%s' is unheard of. Type 'help' for a list of species.",
	"Fie! '%s' is not among the races I've encountered. Type 'help' for a list of species.",
	"Confusion! The species '%s' hides from my gaze. Type 'help' for a list of species.",
	"Oh no! '%s' is not in the wilds I know. Type 'help' for a list of species.",
	"Zounds! '%s' confounds the bravest adventurers. Type 'help' for a list of species.",
	"Egad! The race '%s' hath never crossed my path. Type 'help' for a list of species.",
	"Goodness me! '%s' is as mythical as a unicorn. Type 'help' for a list of species.",
	"Heavens! '%s' is unknown to mortal ken. Type 'help' for a list of species.",
}

var backgroundErrorMessages = []string{
	"Hark! The background '%s' is not in my chronicles. Type 'help' for a list of backgrounds.",
	"Alas! The origin '%s' is lost to the ages. Type 'help' for a list of backgrounds.",
	"By the gods! Such a background '%s' is not in my scrolls. Type 'help' for a list of backgrounds.",
	"Fie! '%s' is not among my tales of heroes. Type 'help' for a list of backgrounds.",
	"Confusion! The background '%s' evades my memory. Type 'help' for a list of backgrounds.",
	"Oh no! '%s' is not in my adventurer's tales. Type 'help' for a list of backgrounds.",
	"Zounds! '%s' baffles the greatest bards. Type 'help' for a list of backgrounds.",
	"Egad! The background '%s' hath never been sung. Type 'help' for a list of backgrounds.",
	"Goodness me! '%s' is as elusive as a shadow. Type 'help' for a list of backgrounds.",
	"Heavens! '%s' is beyond my epic knowledge. Type 'help' for a list of backgrounds.",
}

var classErrorMessages = []string{
	"Hark! The class '%s' is not in my spellbooks. Type 'help' for a list of classes.",
	"Alas! The vocation '%s' is unknown to me. Type 'help' for a list of classes.",
	"By the gods! Such a class '%s' is not in my teachings. Type 'help' for a list of classes.",
	"Fie! '%s' is not among the paths I've walked. Type 'help' for a list of classes.",
	"Confusion! The class '%s' eludes my wisdom. Type 'help' for a list of classes.",
	"Oh no! '%s' is not in my guild's lore. Type 'help' for a list of classes.",
	"Zounds! '%s' confounds the wisest mages. Type 'help' for a list of classes.",
	"Egad! The class '%s' hath never been chosen. Type 'help' for a list of classes.",
	"Goodness me! '%s' is as rare as a dragon's hoard. Type 'help' for a list of classes.",
	"Heavens! '%s' is beyond my arcane grasp. Type 'help' for a list of classes.",
}
