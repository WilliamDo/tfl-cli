package main

import "testing"

func TestFindLine(t *testing.T) {
	filtered, _ := findLine([]Line{
		{ Name: "district" },
		{ Name: "central" },
	}, "district")

	if filtered.Name != "district" {
		t.Fail()
	}
}
