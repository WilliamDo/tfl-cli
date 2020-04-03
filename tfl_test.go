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

	expected := "\u001b[32m\u2713\u001b[0m district...\u001b[32m[very good]\u001b[0m\n"
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

	expected := "\u001b[31m\u2717\u001b[0m district...\u001b[31m[not good]\u001b[0m\n"
	if buffer.String() != expected {
		t.Errorf("want %s; got %s", expected, buffer.String())
	}
}

func TestPrintStatusForMultipleStatus(t *testing.T) {
	buffer := new(bytes.Buffer)
	printStatus(buffer, []Line{
		{ 
			Name: "district", 
			LineStatuses: []LineStatus{
				{ StatusSeverity: 10, StatusSeverityDescription: "very good" },
				{ StatusSeverity: 10, StatusSeverityDescription: "maybe good" },
			},
		},
	})

	expected := "\u001b[32m\u2713\u001b[0m district...\u001b[32m[very good, maybe good]\u001b[0m\n"
	if buffer.String() != expected {
		t.Errorf("want %s; got %s", expected, buffer.String())
	}
}

func TestPrintDepartureBoard(t *testing.T) {
	buffer := new(bytes.Buffer)
	printDepartureBoard(buffer, []Prediction {
		{ Towards: "destination", Direction: "inbound", TimeToStation: 60 },
		{ Towards: "destination", Direction: "inbound", TimeToStation: 90 },
	}, "inbound")

	expected := "destination...1.0\ndestination...1.5\n"
	if buffer.String() != expected {
		t.Errorf("want %s; got %s", expected, buffer.String())
	}
}