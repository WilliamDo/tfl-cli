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
	"text/tabwriter"
)

const (
	INBOUND  string = "inbound"
	OUTBOUND string = "outbound"
)

func main() {

	boardCmd := flag.NewFlagSet("board", flag.ExitOnError)
	boardOutbound := boardCmd.Bool("outbound", false, "outbound")
	boardInbound := boardCmd.Bool("inbound", false, "inbound")
	boardNaptanId := boardCmd.String("naptanId", "", "naptanId")

	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	statusLine := statusCmd.String("line", "", "line")

	stationCmd := flag.NewFlagSet("station", flag.ExitOnError)
	stationList := stationCmd.Bool("list", false, "list")

	out := os.Stdout

	if len(os.Args) < 2 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "status": 
		statusCmd.Parse(os.Args[2:])
		getAndPrintStatus(out, *statusLine)
	case "board": 
		boardCmd.Parse(os.Args[2:])

		if *boardNaptanId == "" {
			fmt.Println("board information needs a station")
			os.Exit(1)
		}

		if *boardOutbound {
			getAndPrintDepartureBoard(out, *boardNaptanId, OUTBOUND)
		}

		if *boardInbound {
			getAndPrintDepartureBoard(out, *boardNaptanId, INBOUND)
		}
	case "station":
		stationCmd.Parse(os.Args[2:])

		if *stationList {
			printNaptanIds(out)
		}
	default:
		fmt.Println("expected 'status' or 'board' subcommands")
		os.Exit(1)
	}
}

type Prediction struct {
	Towards       string
	Direction     string
	TimeToStation float64
}

type Line struct {
	Id           string
	Name         string
	LineStatuses []LineStatus
}

type LineStatus struct {
	StatusSeverity            float64
	StatusSeverityDescription string
}

func getAndPrintStatus(out io.Writer, lineFilter string) {
	resp, err := http.Get("https://api.tfl.gov.uk/Line/Mode/tube/Status")
    if err != nil {
		// handle error
		fmt.Fprintf(out, "error with http\n")
		return
    }
    defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	
	var lines []Line
    jerr := json.Unmarshal(body, &lines)

    if jerr != nil {
		fmt.Fprintf(out, "error with unmarshalling response\n")
		return
	} 

	if lineFilter == "" {
		printStatus(out, lines)
	} else {
		line, err := findLine(lines, lineFilter)
		if err != nil {
			fmt.Fprintf(out, "unrecognised line filter\n")
		} else {
			printStatus(out, []Line{ line })
		}
	}

}

func printStatus(out io.Writer, lines []Line) {

	const padding = 3
	w := tabwriter.NewWriter(out, 0, 0, padding, '.', 0)

	for _, line := range lines {
		concatenatedStatus, statusOk := concatStatus(line.LineStatuses)
		if statusOk {
			fmt.Fprintf(w, "\u001b[32m\u2713\u001b[0m %s\t\u001b[32m[%s]\u001b[0m\n", line.Name, concatenatedStatus)	
		} else {
			fmt.Fprintf(w, "\u001b[31m\u2717\u001b[0m %s\t\u001b[31m[%s]\u001b[0m\n", line.Name, concatenatedStatus)
		}
	}

	w.Flush()
}

func concatStatus(statuses []LineStatus) (string, bool) {
	statusOk := true
	statusDescription := ""
	for index, status := range statuses {
		if status.StatusSeverity == 10 {
	
		} else {
			statusOk = false
		}
		if index > 0 {
			statusDescription = statusDescription + ", "
		}
		statusDescription = statusDescription + status.StatusSeverityDescription
	}

	return statusDescription, statusOk
}

func findLine(lines []Line, lineId string) (Line, error) {
	for _, line := range lines {
		if line.Id == lineId {
			return line, nil
		}
	}

	return Line{}, errors.New("could not find the line")
}

func getAndPrintDepartureBoard(out io.Writer, naptanId string, direction string) {
    resp, err := http.Get("https://api.tfl.gov.uk/StopPoint/" + naptanId + "/Arrivals")
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

	printDepartureBoard(out, predictions, direction)

}

func printDepartureBoard(out io.Writer, predictions []Prediction, direction string) {

	const padding = 3
	w := tabwriter.NewWriter(out, 0, 0, padding, '.', 0)

	for _, train := range predictions {
		if train.Direction == direction {
			fmt.Fprintf(w, "%s\t%.1f\n", train.Towards, train.TimeToStation / 60)
		}
	}

	w.Flush()
}

type StopPointsResponse struct {
	StopPoints []StopPoint
}

type StopPoint struct {
	NaptanId   string
	CommonName string
}

func printNaptanIds(out io.Writer) {
	resp, err := http.Get("https://api.tfl.gov.uk/StopPoint/Mode/tube")
    if err != nil {
		// handle error
		fmt.Printf("error with http")
		return
    }
    defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	
	var stopPointsResponse StopPointsResponse
	jerr := json.Unmarshal(body, &stopPointsResponse)

	if jerr != nil {
		fmt.Printf("error with unmarshalling response")
		return
	} 

	const padding = 3
	w := tabwriter.NewWriter(out, 0, 0, padding, '.', 0)

	for _, stop := range stopPointsResponse.StopPoints {
		fmt.Fprintf(w, "%s\t%s\n", stop.NaptanId, stop.CommonName)
	}

	w.Flush()
}