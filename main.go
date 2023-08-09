package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

type Slot struct {
	Date      string `json:"date"`
	StartTime string `json:"startTime"`
}

type SlotResponse struct {
	Status string `json:"status"`
	Data   []Slot `json:"data"`
}

func main() {
	var desk string
	var persons int
	var productType string

	flag.StringVar(&desk, "d", "", "Desk location:\n AM - Amsterdam\n DH - Den Haag\n ZW - Zwole\n DB - Den Bosch")
	flag.IntVar(&persons, "p", 0, "Number of persons")
	flag.StringVar(&productType, "t", "", "Type of product: \n DOC - Collect residence document \n BIO - Biometrics")
	flag.Parse()

	if persons <= 0 || desk == "" || productType == "" {
		fmt.Printf("Usage: -d <desk_location> -p <number_of_persons> -t <productType> \n")
		flag.PrintDefaults()
		return
	}

	url := fmt.Sprintf("https://oap.ind.nl/oap/api/desks/%s/slots/?productKey=%s&persons=%d", desk, productType, persons)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making the HTTP request:", err)
		return
	}
	defer response.Body.Close()

	var bodyBytes []byte
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	bodyBytes = buf.Bytes()

	// Remove the prefix characters ")]}',"
	bodyBytes = bytes.TrimPrefix(bodyBytes, []byte(")]}',"))

	var slotResponse SlotResponse
	err = json.Unmarshal(bodyBytes, &slotResponse)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return
	}

	for _, slot := range slotResponse.Data {
		fmt.Printf("%s %s\n", slot.Date, slot.StartTime)
	}
}
