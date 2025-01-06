package commands

import (
	"github.com/jake-abed/lore/internal/config"
	"github.com/jake-abed/lore/internal/db"
)

type Command struct {
	Flags       map[string]string
	Name        string
	Description string
	Callback    func(*State) error
}

type State struct {
	Cfg  *config.Config
	Db   *db.Queries
	Args []string
}

func BuildCommands() map[string]Command {
	return map[string]Command{
		"help": {
			Name:        "help",
			Description: "Get information about all available commands.",
			Flags:       nil,
			Callback:    commandHelp,
		},
		"monsters": {
			Name:        "monsters",
			Description: "Get info about monsters and simulate fights.",
			Flags: map[string]string{
				"-i":  "Looks up info about a particular monster by name or id slug.",
				"-f":  "Simulate a fight between two monsters. Name or id slug work.",
				"-va": "View all monsters on the D&D 5e OpenAPI.",
			},
			Callback: commandMonsters,
		},
		"npcs": {
			Name:        "npcs",
			Description: "Add, search, edit, and view info about NPCs.",
			Flags: map[string]string{
				"-v": "Inspect an NPC and view their info.",
				"-a": "Add an NPC to your local database for your campaign.",
				"-e": "Edit an NPC's info by name. Case-insensitive.",
				"-s": "Search your NPCs by name. Returns all possible matches.",
				"-d": "Delete an NPC by name. Case-sensitive.",
			},
			Callback: commandNpcs,
		},
		"places": {
			Name: "places",
			Description: "Add, search, edit, and view worlds, regions, and " +
				"locations. When interacting with, creating, or deleting a " +
				"specific entry, you must pass a 'type' flag such as '--world' " +
				"or '--region'",
			Flags: map[string]string{
				"-a":         "Add a place.",
				"-e":         "Edit a place.",
				"-s":         "Search places by name. Returns all possible matches",
				"-d":         "Delete a place by name. Case-sensitive.",
				"-v":         "Inspect a place and it's information by name.",
				"--world":    "Specify an operation on a world.",
				"--region":   "Specify an operation on a region.",
				"--location": "Specify an operation on a location.",
			},
			Callback: commandPlaces,
		},
	}
}
