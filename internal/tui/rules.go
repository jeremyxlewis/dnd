package tui

var rules = map[string]string{
	"combat": `Combat in D&D 5e follows these steps:
1. Determine surprise. The DM determines if any combatants are surprised.
2. Establish positions. The DM decides where all creatures are located.
3. Roll initiative. Each participant rolls a d20 + Dexterity modifier.
4. Take turns. On your turn: Move, Action, Bonus Action, Reaction.
5. Repeat until combat ends.

Actions: Attack, Cast a Spell, Dash, Disengage, Dodge, Help, Hide, Ready, Search, Use an Object.
Bonus Actions: Off-hand attack, certain spells/features.
Reactions: Opportunity attacks, certain spells/features.`,
	"conditions": `Conditions alter a creature's capabilities in a variety of ways. Common conditions:
- Blinded: Can't see, auto-fail checks requiring sight, attacks have disadvantage, attacks against have advantage.
- Charmed: Can't attack charmer, charmer has advantage on social checks.
- Deafened: Can't hear, auto-fail checks requiring hearing.
- Frightened: Disadvantage on checks/attacks while source in sight, can't willingly move closer.
- Grappled: Speed 0, ends if grappler incapacitated.
- Incapacitated: Can't take actions or reactions.
- Invisible: Can't be seen, attacks have advantage, attacks against have disadvantage.
- Paralyzed: Can't move/speak, auto-fail Str/Dex saves, attacks have advantage, auto-crit if within 5 ft.
- Petrified: Transformed to stone, incapacitated, doesn't age, resistant to all damage.
- Poisoned: Disadvantage on attack rolls and ability checks.
- Prone: Can only crawl, disadvantage on attacks, melee advantage, ranged disadvantage.
- Restrained: Speed 0, disadvantage on attacks, attacks against have advantage, disadvantage on Dex saves.
- Stunned: Incapacitated, can't move, auto-fail Str/Dex saves, attacks against have advantage.
- Unconscious: Incapacitated, unaware, drops everything, falls prone, auto-fail Str/Dex saves, attacks have advantage, auto-crit if within 5 ft.`,
	"ability checks": `Ability checks test a character's innate talent and training. Roll d20 + ability modifier + proficiency bonus (if proficient).
- Strength: Athletics (climbing, swimming, jumping).
- Dexterity: Acrobatics (balance, tumbling), Sleight of Hand (pickpocket), Stealth (hide).
- Constitution: Endurance-related checks.
- Intelligence: Arcana, History, Investigation, Nature, Religion.
- Wisdom: Animal Handling, Insight, Medicine, Perception, Survival.
- Charisma: Deception, Intimidation, Performance, Persuasion.

Advantage/Disadvantage: Roll two d20s, take higher/lower.`,
	"initiative": `At the start of combat, roll initiative: d20 + Dexterity modifier. Higher goes first. Ties broken by Dexterity modifier, then by DM.`,
	"actions": `On your turn, you can take one action. Common actions:
- Attack: Make a melee or ranged attack.
- Cast a Spell: Cast a spell with casting time of 1 action.
- Dash: Gain extra movement equal to your speed.
- Disengage: Movement doesn't provoke opportunity attacks.
- Dodge: Attacks against you have disadvantage, you have advantage on Dex saves.
- Help: Aid an ally's task.
- Hide: Make a Dexterity (Stealth) check.
- Ready: Prepare an action to trigger later.
- Search: Look for something.
- Use an Object: Interact with a second object.`,
}
