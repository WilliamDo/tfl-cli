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

    printDepartureBoard("940GZZLUBKG", OUTBOUND)
	
	printStatus()
	
    fmt.Printf("")
}

type Prediction struct {
	Towards       string
	Direction     string
	TimeToStation float64
}

type Line struct {
	Name string
	LineStatuses []LineStatus
}

type LineStatus struct {
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
		fmt.Println("\u2713 ", line.Name, "\t", line.LineStatuses[0].StatusSeverityDescription)	
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