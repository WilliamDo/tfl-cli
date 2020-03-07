package main

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"
import "sort"

const (
	INBOUND  string = "inbound"
	OUTBOUND string = "outbound"
)

func main() {

    // printDepartureBoard("940GZZLUBKG", OUTBOUND)
	
	printStatus()
	
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

func printStatus() {
	//https://api.tfl.gov.uk/Line/Mode/tube/Status
	resp, err := http.Get("https://api.tfl.gov.uk/Line/Mode/tube/Status")
    if err != nil {
		// handle error
		fmt.Printf("error with http")
		return
    }
    defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	
	var lines []Line
    jerr := json.Unmarshal(body, &lines)

    if jerr != nil {
		fmt.Printf("error with unmarshalling response")
		return
	} 

	for _, line := range lines {
		if line.LineStatuses[0].StatusSeverity == 10 {
			fmt.Println("\u001b[32m\u2713\u001b[0m", line.Name, "\t", "\u001b[32m", line.LineStatuses[0].StatusSeverityDescription, "\u001b[0m")	
		} else {
			fmt.Println("\u001b[31m\u2717\u001b[0m", line.Name, "\t", "\u001b[31m", line.LineStatuses[0].StatusSeverityDescription, "\u001b[0m")
		}
	}
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