# Product Requirements Document: D&D CLI Companion

## Product Overview

**Product Vision:** To provide Dungeons & Dragons players and Dungeon Masters with a fast, easy-to-use command-line tool for looking up rules, rolling dice, and generating game content.

**Target Users:**
- **Primary:** Dungeon Masters (DMs) who need to quickly access information and generate content during a game session.
- **Secondary:** D&D players who want a convenient way to check rules or roll dice without leaving their terminal.

**Business Objectives:**
- Create a useful and fun tool for the D&D community.
- Serve as a practical project for learning and practicing CLI/TUI application development in Go.
- Release as an open-source project to encourage community contributions.

**Success Metrics:**
- Positive feedback and adoption from the D&D and open-source communities.
- Active use of the tool in personal D&D campaigns.
- Growth in community contributions (e.g., forks, pull requests) if open-sourced.

## User Personas

### Persona 1: "DM Dave"
- **Demographics:** 32-year-old software developer, high technical proficiency.
- **Goals:** Keep his D&D sessions flowing smoothly without getting bogged down in rulebooks. Needs to improvise NPCs and story elements on the fly.
- **Pain Points:** Flipping through the physical Player's Handbook (PHB) is slow and disruptive. Coming up with unique NPC names and simple backstories under pressure is challenging.
- **User Journey:** During a game, a player asks about the "Grappled" condition. Dave types `dnd lookup grappled` and gets the rules instantly. Later, the players enter a tavern, and Dave types `dnd npc generate` to get a quick NPC to tend the bar.

### Persona 2: "Player Penny"
- **Demographics:** 26-year-old student, moderate technical proficiency.
- **Goals:** Quickly check a spell's description or roll dice without having to find her physical dice or use a clunky web app.
- **Pain Points:** Forgetting the exact wording or mechanics of a spell. Not having the right type or number of dice on hand.
- **User Journey:** Penny is in a heated battle and wants to cast "Fireball." She types `dnd spell fireball` to double-check the area of effect. She then types `dnd roll 8d6` to roll the damage.

## Feature Requirements

| Feature | Description | User Stories | Priority | Acceptance Criteria | Dependencies |
|---|---|---|---|---|---|
| **Character Management** | Comprehensive tools for creating, managing, and leveling up D&D characters. | As a user, I want to create a character step-by-step so I can quickly get started. As a user, I want to save and load character sheets so I can track my progress. As a user, I want to level up my character so their stats and abilities are automatically updated. | Must | Guided character creation process. Ability to save/load character data. Automatic calculation of HP, proficiencies, and spell slots on level-up. | Local storage for character data. |
| **Information Lookup** | A robust database for quick access to D&D rules, spells, monsters, and items. | As a user, I want to look up specific rules so I can understand game mechanics. As a user, I want to search for spell descriptions so I know their effects. As a DM, I want to find monster stat blocks so I can run encounters. As a user, I want to browse items and equipment so I can manage my inventory. | Must | `dnd rule <topic>` returns relevant rule text. `dnd spell <name>` returns full spell details. `dnd monster <name>` returns full stat block. `dnd item <name>` returns item description. | D&D 5e SRD data from `nick-aschenbach/dnd-data`. |
| **Combat & Encounter Management** | Tools to streamline combat encounters for DMs and players. | As a user, I want to roll dice with modifiers and advantage/disadvantage so I can resolve actions. As a DM, I want to track initiative so I can manage turn order. As a DM, I want to track health and conditions so I can keep combat organized. As a DM, I want to generate balanced encounters so I can challenge my players appropriately. | Must | `dnd roll <dice_notation>` supports modifiers, advantage/disadvantage. `dnd init add/next/list` manages initiative. `dnd combat hp <char> <+/-amount>` updates health. `dnd encounter generate <party_level> <difficulty>` suggests monsters. | Monster data from `nick-aschenbach/dnd-data`. |
| **Utilities** | Helpful supplementary features for D&D sessions. | As a DM, I want to generate random names/loot/encounters so I can improvise quickly. As a user, I want to save and load the state of my game so I can resume later. | Should | `dnd generate name/loot/encounter`. `dnd session save/load`. | Data for random generation. |

## User Flows

### Flow 1: Rolling Dice for an Attack
1. User types `dnd roll 1d20+7`.
2. The tool calculates the sum of one 20-sided die and adds 7.
3. The tool prints the result to the console, e.g., `15`.

### Flow 2: Looking up a Spell
1. User types `dnd spell "eldritch blast"`.
2. The tool searches its data source for the "Eldritch Blast" spell.
3. The tool prints the spell's description, range, components, and other details in a clean, readable format.
4. **Error state:** If the spell is not found, the tool prints a D&D-themed error message.

### Flow 3: Creating a New Character
1. User types `dnd char create`.
2. The tool prompts the user through ability score generation (e.g., roll 4d6 drop lowest, assign to stats).
3. The tool prompts for species, class, background, and other character details.
4. The tool saves the complete character sheet to a local file.

### Flow 4: Tracking Combat Initiative
1. DM types `dnd init add "Goblin 1" 12`.
2. DM types `dnd init add "Player A" 18`.
3. DM types `dnd init list` to see the current order.
4. DM types `dnd init next` to advance to the next turn.

## User Flows

### Flow 1: Rolling Dice for an Attack
1. User types `dnd roll 1d20+7`.
2. The tool calculates the sum of one 20-sided die and adds 7.
3. The tool prints the result to the console, e.g., `15`.

### Flow 2: Looking up a Spell
1. User types `dnd spell "eldritch blast"`.
2. The tool searches its data source for the "Eldritch Blast" spell.
3. The tool prints the spell's description, range, components, and other details in a clean, readable format.
4. **Error state:** If the spell is not found, the tool prints a D&D-themed error message.

### Flow 3: Creating a New Character
1. User types `dnd char create`.
2. The tool prompts the user through ability score generation (e.g., roll 4d6 drop lowest, assign to stats).
3. The tool prompts for species, class, background, and other character details.
4. The tool saves the complete character sheet to a local file.

### Flow 4: Tracking Combat Initiative
1. DM types `dnd init add "Goblin 1" 12`.
2. DM types `dnd init add "Player A" 18`.
3. DM types `dnd init list` to see the current order.
4. DM types `dnd init next` to advance to the next turn.

## TUI User Flow & UX Design

This section outlines the proposed user experience and interaction flows for the Terminal User Interface (TUI) component of the D&D CLI Companion. The TUI aims to provide a more interactive and visually rich experience for certain features, while maintaining the efficiency of a command-line interface.

### High-Level Flow

1.  **Main Command Input Screen:**
    *   This will be the default view when the TUI is launched.
    *   Users will type commands (e.g., `roll 2d6`, `spell fireball`) into a single, prominent text input field.
    *   The output of executed commands will appear in a dedicated area above the input field.
    *   **Discoverability:** A persistent hint or a dedicated `help` command (`help` or `?`) will list available commands and features.
    *   **Styling:** `lipgloss` will be used to clearly separate and style the input, output, and status areas, enhancing readability and visual appeal.

2.  **Feature-Specific Modes/Screens:**
    *   For more complex, interactive features (e.g., Character Management, Initiative Tracking), the TUI will transition into dedicated "modes" or "screens."
    *   These modes will utilize `charmbracelet/bubbles` components (e.g., lists, forms, tables) to provide richer, guided interactions.
    *   **Navigation:** Clear mechanisms will be provided to exit a mode (e.g., `Esc`, `Ctrl+C`, or a specific command like `back` or `exit`) and return to the Main Command Input Screen.

### Detailed Feature UX Suggestions

*   **Information Lookup (`lookup <topic>`, `spell <name>`, `monster <name>`, `item <name>`):**
    *   **Flow:** User types command -> TUI executes the corresponding CLI command -> displays formatted output.
    *   **UX:** Output will be paginated if lengthy (using `bubbles/viewport`), syntax-highlighted (if applicable), and clearly structured (e.g., bold headings, bullet points, distinct sections).
    *   **Enhancement:** Autocompletion for spell, monster, or item names within the input field (leveraging `bubbles/textinput` with a custom completer) will improve efficiency.
    *   **Fuzzy Find Enhancement:** Implement a fuzzy finder for commands like `spell`, `monster`, and `item` lookup. If a user types a partial name or a command without a specific argument (e.g., `spell`), the TUI will present an interactive fuzzy search interface. Users can type to filter a list of available items (spells, monsters, etc.) and select the desired entry using arrow keys and Enter. This significantly improves discoverability and usability when exact names are unknown.

*   **Dice Rolling (`roll <notation>`):**
    *   **Flow:** User types `roll 2d6+3` -> TUI executes the `roll` command -> displays the result.
    *   **UX:** Clear display of individual dice rolls and the final total. Visual flair (e.g., a subtle animation or color change) could be considered for a more engaging experience.

*   **Character Management (`char create`, `char load`, `char view`, `char levelup`):**
    *   **Flow:** User types `char create` -> TUI transitions to a "Character Creation Mode."
    *   **UX (Character Creation Mode):**
        *   Utilize `bubbles/form` or a guided series of `bubbles/textinput` and `bubbles/select` components.
        *   Guide the user step-by-step through character creation (name, species, class, ability scores, background, etc.).
        *   Display a summary of the character sheet as it is being built.
        *   Provide clear options for saving and loading character data.
    *   **UX (Character View Mode):** Display a well-formatted, potentially scrollable, character sheet with key stats, abilities, and inventory.

*   **Combat & Encounter Management (`init track`, `combat hp`, `encounter generate`):**
    *   **Flow:** User types `init track` -> TUI transitions to an "Initiative Tracker Mode."
    *   **UX (Initiative Tracker Mode):**
        *   Use `bubbles/table` or `bubbles/list` to display the initiative order.
        *   Commands within this mode (e.g., `add <name> <init>`, `next`, `remove <name>`, `clear`) will be context-sensitive.
        *   Clearly highlight the current turn.
    *   **UX (HP/Condition Tracker):** Similar to the initiative tracker, but focused on tracking hit points, temporary hit points, and various conditions (e.g., poisoned, prone, grappled) for combatants.

*   **Utilities (`generate name`, `generate loot`, `session save`, `session load`):**
    *   **Flow:** User types command -> TUI executes the utility command -> displays generated content.
    *   **UX:** Clear, concise output. For generators, a list of options or a single generated item will be presented.

### Overall UX Principles

*   **Consistency:** Maintain consistent command syntax, navigation patterns, and visual styling across all TUI modes and features.
*   **Feedback:** Provide immediate and clear feedback for all user actions, indicating success, failure, or ongoing processes.
*   **Error Handling:** Integrate the D&D-themed error messages (as per the PRD) into the TUI's output.
*   **Responsiveness:** Ensure the TUI feels snappy and responsive to user input, minimizing perceived lag.
*   **Theming:** Allow for customizable color themes (via `lipgloss`) to cater to user preferences.

### Immediate Next Steps for TUI Implementation

1.  **Implement a simple "help" command:** This will be a crucial first step to demonstrate command processing within the TUI and provide initial discoverability for users.
2.  **Refine `lastOutput` display:** Ensure it handles multi-line output gracefully, potentially clearing after a few commands or becoming scrollable for extensive results.
3.  **Consider a "mode" manager:** Develop a higher-level `tea.Model` that can effectively switch between different sub-models (e.g., `MainInputModel`, `CharCreationModel`, `InitiativeModel`) to manage complex feature interactions.

## Non-Functional Requirements

### Performance
- **Response Time:** All commands should execute in under 500ms.

### Security
- **Data Protection:** Character sheet data and saved session data should be stored locally in the user's home directory in a non-encrypted, human-readable format (e.g., JSON). No sensitive data is handled.

### Compatibility
- **OS:** Should be cross-platform and run on Linux, macOS, and Windows.
- **Terminal:** Should work in any standard terminal emulator.

### Accessibility
- **Compliance Level:** Standard command-line interface accessibility. Output should be clear and well-formatted.

### Error Handling
- **D&D-Themed Feedback:** When an error or an invalid query occurs, the tool shall respond with a message befitting a seasoned adventurer's plea for clarity, such as: "Hark! Thy query, good sir or madam, doth bewilder my arcane senses. Pray tell, couldst thou rephrase thy plea, for its meaning doth elude my understanding?" or similar evocative phrasing.

## Technical Specifications

### Application
- **Technology Stack:** Go, utilizing `charmbracelet/bubbletea` for any interactive Terminal User Interface (TUI) components. CLI command parsing will use a library such as `cobra` or `urfave/cli`. Data parsing and serialization will leverage Go's standard library `encoding/json` or a dedicated library like `json-iterator/go` for performance.
- **Data Storage:** Game data (rules, spells, monsters, items) will be sourced from `https://github.com/nick-aschenbach/dnd-data` (JSON files) and loaded into memory or a local cache for efficient lookup. User data (character sheets, saved sessions) will be stored as individual JSON files in the user's home directory.

### Infrastructure
- **Hosting:** The code will be hosted on a public Git repository (e.g., GitHub).
- **CI/CD:** A simple CI pipeline will be set up to run tests on each push.
- **Deployment/Distribution:** Distribution will primarily be via GitHub Releases, providing pre-compiled binaries for various operating systems.

## Release Planning

### MVP (v1.0)
- **Features:**
  - Dice Roller (with modifiers, advantage/disadvantage).
  - Information Lookup (for core rules, conditions, all spells, and a subset of monsters/items from `nick-aschenbach/dnd-data`).
- **Timeline:** 4-6 weeks.
- **Success Criteria:** The tool can successfully roll dice with various options and provide comprehensive lookup for spells and basic rules from the specified SRD data source.

### Future Releases
- **v1.1:** Expand Information Lookup to include all SRD monsters and items. Implement basic NPC Generator.
- **v1.2:** Implement Character Management (creation, saving, loading, basic leveling).
- **v2.0:** Implement Combat & Encounter Management (Initiative Tracker, Health/Status Tracker, basic Encounter Builder).
- **v2.1+:** Advanced features like Saved Sessions, more robust Random Generators, custom content integration, and potential for more elaborate TUI elements using `bubbletea`.

## Open Questions & Assumptions

- **Question 1:** What specific Go libraries will be chosen for CLI parsing (`cobra` vs. `urfave/cli`) and data handling?
- **Assumption 1:** Users will have Go installed on their systems (or use provided binaries).
- **Assumption 2:** Users are comfortable working with a command-line interface.
- **Assumption 3:** The `nick-aschenbach/dnd-data` repository provides sufficient and accurate data for the initial versions of the lookup features.

## Appendix

### Legal & Licensing
- **Project License:** The project will be released under an appropriate open-source license (e.g., MIT License).
- **SRD Compliance:** All usage of D&D 5e SRD content will strictly adhere to the Open Game License (OGL) requirements, including proper attribution and inclusion of the OGL text within the project.

### Testing Strategy
- **Unit Tests:** Comprehensive unit tests will be written for all core logic, including dice rolling, data parsing, and command handling, using Go's built-in testing framework.
- **Integration Tests:** Tests will verify the interaction between different modules and the correct execution of CLI commands.
- **Acceptance Tests:** High-level tests will ensure that key user stories are met from an end-to-end perspective.

### Future Considerations (Beyond Roadmap)
- Integration with external D&D platforms (if APIs become available).
- Support for homebrew content creation and management.
- More advanced TUI interfaces for specific features (e.g., character sheet editor, initiative tracker).
- Multi-user support for shared game sessions (complex, long-term goal).
