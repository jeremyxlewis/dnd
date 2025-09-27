package tui

// Default dimensions
const (
	DefaultWidth          = 80
	DefaultHeight         = 24
	ViewportWidthPadding  = 2
	ViewportHeightPadding = 3
	ListHeightPadding     = 4
	TextInputCharLimit    = 156
	TextInputWidth        = 40
	HistoryLimit          = 100
)

// Step constants for charCreateModel
const (
	StepName = iota
	StepAlignment
	StepPlayer
	StepLevel
	StepScoreMethod
	StepScores
	StepSpecies
	StepSpeciesInfo
	StepClass
	StepClassInfo
	StepBackground
	StepBackgroundInfo
	StepProficiencies
	StepEquipment
	StepSpellcasting
	StepConfirm
)

// Input modes for initiativeTracker
const (
	InputModeAddName = "add_name"
	InputModeAddInit = "add_init"
)
