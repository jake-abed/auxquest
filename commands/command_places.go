package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/jake-abed/lore/internal/db"
)

func commandPlaces(s *State) error {
	args := s.Args[1:]

	if len(args) < 2 {
		return fmt.Errorf("Places command requires at least two arguments!")
	}

	var typeFlag string

	for _, arg := range args {
		if isPlaceTypeFlag(arg) {
			typeFlag = arg
		}
	}

	flag, flagArg := parseFlagArg(args)

	if flag != "-s" && typeFlag == "" {
		return fmt.Errorf("Flag %s requires a place type flag as well.", flag)
	} else if flag == "-s" && flagArg != "" {
		fmt.Println("Add `search fn!")
		return nil
	}

	switch flag {
	case "-a":
		place, err := addPlace(s, typeFlag)
		if err != nil {
			fmt.Println("Uh oh! Lore errored out while adding this place: ")
			return err
		}
		printPlace(place)
		return nil
	case "-v":
		place, err := getPlaceByName(s, typeFlag, strings.ToLower(flagArg))
		if err != nil {
			fmt.Printf("Hmm... Lore couldn't find %s. Here's the error: \n", flagArg)
			return err
		}
		world, _ := place.(*db.World)
		printWorld(world)
		return nil
	case "-e":
		fmt.Println("Add `edit` fn!")
		return nil
	case "-d":
		fmt.Println("Add `delete` fn!")
		return nil
	default:
		fmt.Println("Help!")
		return nil
	}
}

// Place Printers

func printPlace(p db.Place) {
	switch p.(type) {
	case *db.World:
		world := p.(*db.World)
		printWorld(world)
	case *db.Area:
		area := p.(*db.Area)
		printArea(area)
	default:
		fmt.Printf("Lore has no such place type as %T\n", p)
	}
}

func printWorld(w *db.World) {
	headerMsg := fmt.Sprintf("World: %-16s Id: %-2d", w.Name, w.Id)
	printHeader(headerMsg)
	fmt.Println(bold.Render("Description: ") + w.Desc)
}

func printArea(r *db.Area) {
	headerMsg := fmt.Sprintf("Area: %-16s Id: %-2d", r.Name, r.Id)
	printHeader(headerMsg)
	fmt.Println(bold.Render("Area Type: ") +r .Type)
	fmt.Println(bold.Render("Description: ") + r.Desc)
	fmt.Println(bold.Render("Belongs to World Id: ") +
		fmt.Sprintf("%d", r.WorldId))
}

func addPlace(s *State, typeFlag string) (db.Place, error) {
	switch typeFlag {
	case "--world":
		world := worldForm(db.World{})
		worldParams := db.WorldParams{Name: world.Name, Desc: world.Desc}

		newWorld, err := s.Db.AddWorld(context.Background(), &worldParams)
		if err != nil {
			return &db.World{}, err
		}

		return newWorld, nil
	case "--area":
		area := areaForm(s, db.Area{})
		areaParams := &db.AreaParams{
			Name:    area.Name,
			Type:    area.Type,
			Desc:    area.Desc,
			WorldId: area.WorldId,
		}

		newArea, err := s.Db.AddArea(context.Background(), areaParams)
		if err != nil {
			return &db.Area{}, err
		}

		return newArea, err
	default:
		return nil, fmt.Errorf("%s is not a valid typeflag", typeFlag)
	}
}

func getPlaceByName(s *State, typeFlag string, arg string) (db.Place, error) {
	switch typeFlag {
	case "--world":
		world, err := s.Db.GetWorldByName(context.Background(), arg)
		if err != nil {
			return nil, err
		}
		return world, nil
	default:
		return nil, nil
	}
}

// Form Functions

func worldForm(world db.World) db.World {
	huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("World Name: ").
				Value(&world.Name),
			huh.NewText().
				Title("World Description: ").
				Value(&world.Desc),
		),
	).WithTheme(huh.ThemeBase16()).Run()

	return world
}

func areaForm(s *State, area db.Area) db.Area {
	worlds, _ := s.Db.GetXWorlds(context.Background(), 10, 0)
	huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Area Name: ").
				Value(&area.Name),
			huh.NewInput().
				Title("Area Type: ").
				Value(&area.Type),
			huh.NewText().
				Title("Area Description: ").
				Value(&area.Desc),
		),
		newWorldSelectGroup(worlds, area.Name, &area.WorldId),
	).WithTheme(huh.ThemeBase16()).Run()

	return area
}

func newWorldSelectGroup(worlds []*db.World, name string, val *int) *huh.Group {
	options := []huh.Option[int]{}

	for _, world := range worlds {
		option := huh.NewOption(fmt.Sprintf("%d - %s", world.Id, world.Name),
			world.Id,
		)

		options = append(options, option)
	}

	return huh.NewGroup(
		huh.NewSelect[int]().
			Title(fmt.Sprintf("Which world does %s belong to?", name)).
			Options(options...).
			Value(val),
	)
}

// Flag helper functions

func isPlaceTypeFlag(flag string) bool {
	return flag == "--world" || flag == "--area" || flag == "-location"
}

func isCommandFlag(flag string) bool {
	return flag == "-a" || flag == "-v" || flag == "-e" ||
		flag == "-d" || flag == "-s"
}

func parseFlagArg(args []string) (string, string) {
	for i, arg := range args {
		if isCommandFlag(arg) && (1+i) < len(args) {
			return arg, args[i+1]
		} else if isCommandFlag(arg) {
			return arg, ""
		}
	}

	return "", ""
}
