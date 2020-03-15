package main

import (
	"bytes"
	"testing"
)

func TestFindLine(t *testing.T) {
	filtered, _ := findLine([]Line{
		{ Name: "district" },
		{ Name: "central" },
	}, "district")

	if filtered.Name != "district" {
		t.Fail()
	}
}

func TestPrintStatusForEmptyStatus(t *testing.T) {
	buffer := new(bytes.Buffer)
	printStatus(buffer, []Line{})

	if buffer.String() != "" {
		t.Fail()
	}
}

func TestPrintStatusForGoodStatus(t *testing.T) {
	buffer := new(bytes.Buffer)
	printStatus(buffer, []Line{
		{ 
			Name: "district", 
			LineStatuses: []LineStatus{
				{ StatusSeverity: 10 },
			},
		},
	})

	expected := "\u001b[32m\u2713\u001b[0m district\t\u001b[32m\u001b[0m\n"
	if buffer.String() != expected {
		t.Errorf("want %s; got %s", expected, buffer.String())
	}
}

func TestPrintStatusForDelayStatus(t *testing.T) {
	buffer := new(bytes.Buffer)
	printStatus(buffer, []Line{
		{ 
			Name: "district", 
			LineStatuses: []LineStatus{
				{ StatusSeverity: 20 },
			},
		},
	})

	expected := "\u001b[31m\u2717\u001b[0m district\t\u001b[31m\u001b[0m\n"
	if buffer.String() != expected {
		t.Errorf("want %s; got %s", expected, buffer.String())
	}
}