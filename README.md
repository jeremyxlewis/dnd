# D&D CLI Companion

A command-line companion for Dungeons & Dragons players and Dungeon Masters, built with Go and Cobra.

## Features

*   **Dice Roller:** Roll dice using standard D&D notation (e.g., `2d6`, `1d20+5`), with support for advantage and disadvantage.
*   **Information Lookup:** Quickly look up details for spells, monsters, and items from the D&D 5e SRD.
*   **Character Management:** Full D&D character sheet management with Player's Handbook coverage - create, view, level up, track HP/spells/inventory, all 12 classes, 40+ races, 13 backgrounds, complete mechanics, and runtime play tracking.
*   **NPC Generator:** Generate random NPC names, species, and backgrounds.
*   **Interactive TUI:** Full-screen terminal interface for browsing and searching content with fuzzy matching, including guided character creation following D&D rules.

## Installation

### Prerequisites

*   Go (version 1.16 or higher)
*   Git

### Steps

1.  **Clone the repository:**
    ```bash
    git clone --recurse-submodules https://github.com/jeremyxlewis/dnd.git
    cd dnd-cli
    ```

2.  **The D&D data is already included** in the `data/` directory as a git submodule.

3.  **Build the application:**
    ```bash
    go mod tidy
    go build -o dnd .
    ```

## Usage

Run the `dnd` executable from the project root. Use `./dnd` if it's not in your PATH.

### Optional: Install Globally

To run `dnd` from anywhere without `./`:

- Move to a PATH directory: `sudo mv dnd /usr/local/bin/`
- Or add to your PATH: Add `export PATH="$PATH:/path/to/dnd-cli"` to `~/.bashrc` (replace with actual path).

Then use `dnd` or `dnd tui` directly.

### Dice Roller

Roll dice with standard notation:

```bash
dnd roll 1d20
dnd roll 2d6+3
```

Roll with advantage:

```bash
dnd roll 1d20 --advantage
dnd roll 1d20 -a
```

Roll with disadvantage:

```bash
dnd roll 1d20 --disadvantage
dnd roll 1d20 -d
```

### Spell Lookup

Look up a spell by name:

```bash
dnd spell "Fireball"
dnd spell "eldritch blast"
```

### Monster Lookup

Look up a monster by name:

```bash
dnd monster "Goblin"
dnd monster "Ancient Red Dragon"
```

### Item Lookup

Look up an item by name:

```bash
dnd item "Potion of Healing"
dnd item "Longsword"
```

### NPC Generator

Generate a random NPC:

```bash
dnd npc generate
# or simply
dnd npc
```

### Character Management

The D&D CLI provides comprehensive character sheet management for full D&D 5e gameplay, including creation, progression, and runtime tracking.

#### Creating a Character
Create a new character interactively (prompts for species, class, and background with full PHB options - 40+ races, 12 classes, 13 backgrounds):

```bash
dnd char create "Eldrin"
```

This applies all racial traits, class proficiencies/features, background equipment, and initializes HP/spell slots.

#### Viewing Character Sheets
View a character's full sheet (ability scores with modifiers, proficiencies, features, spell slots, equipment, current HP, conditions, etc.):

```bash
dnd char view "Eldrin"
```

Displays complete character details, including current HP (e.g., `HP: 25/30 (+5 temp)`), used spell slots, and inventory.

#### Leveling Up
Level up a character (applies class features, HP increases, spell slots, subclass choices, etc.):

```bash
dnd char levelup "Eldrin"
```

Automatically applies level-appropriate mechanics and updates the character.

#### Managing HP During Play
Track hit points in combat and exploration:

```bash
dnd char hp "Eldrin" damage 10  # Take 10 damage
dnd char hp "Eldrin" heal 5     # Heal 5 HP
dnd char hp "Eldrin" set 20     # Set HP to exact value
```

#### Managing Spell Slots
Track spell usage for spellcasters:

```bash
dnd char spells "Eldrin" use 1 1     # Use 1 first-level slot
dnd char spells "Eldrin" restore 1 1 # Restore 1 first-level slot
```

#### Managing Inventory
Add or remove items from equipment:

```bash
dnd char inventory "Eldrin" add "Potion of Healing"
dnd char inventory "Eldrin" remove "Rusty Sword"
```

#### Complete Character Management Guide

1. **Create Your Character:**
   ```bash
   dnd char create "MyHero"
   # Follow prompts for race, class, background
   ```

2. **Review Your Sheet:**
   ```bash
   dnd char view "MyHero"
   # Check stats, proficiencies, features
   ```

3. **Play the Game:**
   - **Combat:** `dnd char hp "MyHero" damage 15`
   - **Magic:** `dnd char spells "MyHero" use 1 1`
   - **Loot:** `dnd char inventory "MyHero" add "Magic Sword"`

4. **Level Up:**
   ```bash
   dnd char levelup "MyHero"
   # Gains new features, HP, spell slots
   ```

5. **Rest and Recover:**
   ```bash
   dnd char spells "MyHero" restore 1 2  # Long rest recovery
   dnd char hp "MyHero" heal 10          # Healing
   ```

Characters are saved as JSON files in `~/.dnd-cli/` with complete D&D 5e mechanics for all PHB races, classes, and backgrounds, including runtime tracking for HP, spells, inventory, conditions, and inspiration.

## TUI (Terminal User Interface)

Launch the interactive TUI for a more immersive experience:

```bash
dnd tui
```

### Features in TUI

- **Browse All PHB Content:** Browse through all the Player's Handbook content, including spells, monsters, items, and more.
- **Fuzzy Search:** Type `spell`, `monster`, or `item` to browse lists with real-time filtering. Start typing to narrow down results.
- **Character Creation:** Guided step-by-step character creation following D&D 5e rules - select name, alignment, species (with racial traits preview), class (with features), background, ability scores (Standard Array, Roll, or Point Buy), proficiencies, equipment, and spellcasting.
- **Navigation:** Use ↑↓ or jk keys to scroll lists, Enter to select, / to search within lists, Esc to go back.
- **Scrolling:** Use ↑/↓ keys to scroll through long content in the output area.
- **Help:** Type `help` or `?` for available commands.
- **Quitting:** Press Ctrl+C to exit from any screen. Press Esc to exit main or go back in sub-menus.

### Example TUI Workflow

1. Run `dnd tui`
2. Type `spell` and press Enter to browse spells.
3. Type "fire" to filter to fireball-related spells.
4. Use arrows to select, Enter to view details.
5. Press Esc to return to main prompt.
6. Type `roll 1d20` for dice rolls.
7. Select "Create Character" from the main menu to start guided creation.

The TUI provides themed error messages and a clean, scrollable interface for all CLI features.

## Contributing

### Future Improvements

- [x] Replace deprecated `ioutil` functions with `os` equivalents in character handling
- [x] Optimize random seeding to avoid reseeding on every operation
- [x] Implement full character creation and management features
- [x] Expand test coverage for core functionality
- [x] Refactor error message generation to reduce code duplication

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
