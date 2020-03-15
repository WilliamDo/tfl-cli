package main

import (
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"encoding/json"
	"sort"
	"os"
	"flag"
	"errors"
)

const (
	INBOUND  string = "inbound"
	OUTBOUND string = "outbound"
)

// upney: 940GZZLUUPY
// barking: 940GZZLUBKG

var stations = map[string]string {
	"upney":   "940GZZLUUPY",
	"barking": "940GZZLUBKG",
}



func main() {

	boardCmd := flag.NewFlagSet("board", flag.ExitOnError)
	boardOutbound := boardCmd.Bool("outbound", false, "outbound")
	boardInbound := boardCmd.Bool("inbound", false, "inbound")
	boardStation := boardCmd.String("station", "", "station")

	out := os.Stdout

	if len(os.Args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "status": 
		getAndPrintStatus(out)
	case "board": 
		boardCmd.Parse(os.Args[2:])

		if *boardStation == "" {
			fmt.Println("board information needs a station")
			os.Exit(1)
		}

		if *boardOutbound {
			printDepartureBoard(stations[*boardStation], OUTBOUND)
		}

		if *boardInbound {
			printDepartureBoard(stations[*boardStation], INBOUND)
		}
	default:
		fmt.Println("expected 'status' or 'board' subcommands")
		os.Exit(1)
	}

    fmt.Printf("")
}

type Prediction struct {
	Towards       string
	Direction     string
	TimeToStation float64
}

type Line struct {
	Name         string
	LineStatuses []LineStatus
}

type LineStatus struct {
	StatusSeverity            float64
	StatusSeverityDescription string
}

func getAndPrintStatus(out io.Writer) {
	resp, err := http.Get("https://api.tfl.gov.uk/Line/Mode/tube/Status")
    if err != nil {
		// handle error
		fmt.Fprintf(out, "error with http")
		return
    }
    defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	
	var lines []Line
    jerr := json.Unmarshal(body, &lines)

    if jerr != nil {
		fmt.Fprintf(out, "error with unmarshalling response")
		return
	} 

	printStatus(out, lines)

}

func printStatus(out io.Writer, lines []Line) {
	for _, line := range lines {
		if line.LineStatuses[0].StatusSeverity == 10 {
			fmt.Fprintf(out, "\u001b[32m\u2713\u001b[0m %s\t\u001b[32m%s\u001b[0m\n", line.Name, line.LineStatuses[0].StatusSeverityDescription)	
		} else {
			fmt.Fprintf(out, "\u001b[31m\u2717\u001b[0m %s\t\u001b[31m%s\u001b[0m\n", line.Name, line.LineStatuses[0].StatusSeverityDescription)
		}
	}
}

func findLine(lines []Line, name string) (Line, error) {
	for _, line := range lines {
		if line.Name == name {
			return line, nil
		}
	}

	return Line{}, errors.New("could not find the line")
}

func printDepartureBoard(naptanId string, direction string) {
    resp, err := http.Get("https://api.tfl.gov.uk/Line/district/Arrivals/" + naptanId)
    if err != nil {
		// handle error
		fmt.Printf("error with http")
		return
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    var predictions []Prediction
    jerr := json.Unmarshal(body, &predictions)

    if jerr != nil {
		fmt.Printf("error with unmarshalling response")
		return
	} 

	sort.Slice(predictions[:], func(i, j int) bool {
		return predictions[i].TimeToStation < predictions[j].TimeToStation
	})

	for _, train := range predictions {
		if train.Direction == direction {
			fmt.Println(train.Towards, "\t", train.TimeToStation / 60)
		}
	}
}

type StopPoint struct {
	NaptanId   string
	CommonName string
}

func printNaptanIds(lineId string) {
	resp, err := http.Get("https://api.tfl.gov.uk/Line/" + lineId + "/StopPoints")
    if err != nil {
		// handle error
		fmt.Printf("error with http")
		return
    }
    defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	
	var stopPoints []StopPoint
	jerr := json.Unmarshal(body, &stopPoints)

	if jerr != nil {
		fmt.Printf("error with unmarshalling response")
		return
	} 

	for _, stop := range stopPoints {
		fmt.Println(stop.NaptanId, stop.CommonName)
	}
}