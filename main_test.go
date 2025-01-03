package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestMainDefault(t *testing.T) {
	origOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Args = []string{"auxquest"}
	os.Stdout = w
	main()
	os.Stdout = origOut
	w.Close()
	out, _ := io.ReadAll(r)
	readableOutput := string(out)
	expected := []string{
		"Welcome to AuxQuest!",
		"monsters     <==>",
		"*** -i",
		"- View all monsters on the D&D 5e OpenAPI.",
		"Get information about all available commands.",
	}
	for _, phrase := range expected {
		if !strings.Contains(readableOutput, phrase) {
			t.Errorf("Expected StdOut to contain %s", phrase)
		}
	}
}

func TestMainMonsters(t *testing.T) {

	origOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Args = []string{"auxquest", "monsters"}
	os.Stdout = w
	main()
	os.Stdout = origOut
	w.Close()
	out, _ := io.ReadAll(r)
	readableOutput := string(out)
	expected := []string{
		"AuxQuest Monsters Help",
		"Monsters subcommands information",
		"*** monsters -i <monster-name>",
	}
	for _, phrase := range expected {
		if !strings.Contains(readableOutput, phrase) {
			t.Errorf("Expected StdOut to contain %s", phrase)
		}
	}
}
