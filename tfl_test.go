package main

import (
	"bytes"
	"testing"
)

func TestFindLine(t *testing.T) {
	filtered, _ := findLine([]Line{
		{ Id: "district" },
		{ Id: "central" },
	}, "district")

	if filtered.Id != "district" {
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
				{ StatusSeverity: 10, StatusSeverityDescription: "very good" },
			},
		},
	})

	expected := "\u001b[32m\u2713\u001b[0m district...\u001b[32mvery good\u001b[0m\n"
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
				{ StatusSeverity: 20, StatusSeverityDescription: "not good" },
			},
		},
	})

	expected := "\u001b[31m\u2717\u001b[0m district...\u001b[31mnot good\u001b[0m\n"
	if buffer.String() != expected {
		t.Errorf("want %s; got %s", expected, buffer.String())
	}
}

func TestPrintDepartureBoard(t *testing.T) {
	buffer := new(bytes.Buffer)
	printDepartureBoard([]Prediction {
		{ Towards: "destination", Direction: "inbound", TimeToStation: 0 },
	}, "inbound")

	expected := "destination\t0"
	if buffer.String() != expected {
		t.Errorf("want %s; got %s", expected, buffer.String())
	}
}