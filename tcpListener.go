package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

type Data1 struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func detectFormat1(data []byte) (string, error) {

	var jsonObj interface{}
	if err := json.Unmarshal(data, &jsonObj); err == nil {
		return "JSON", nil
	}

	// Try unmarshaling as XML
	var xmlObj interface{}
	if err := xml.Unmarshal(data, &xmlObj); err == nil {
		return "XML", nil
	}

	return "", errors.New("unknown format")
}
func handleConnection(conn net.Conn) {
	defer conn.Close()

	dataRaw, err := fetchDataFromAPI()
	if err != nil {
		fmt.Println("Error fetching data from API:", err)
		return
	}

	id, err := detectFormat(dataRaw)
	if id == "JSON" {
		var data Data
		if err := json.Unmarshal([]byte(dataRaw), &data); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		// Convert to XML
		xmlData, err := xml.MarshalIndent(data, "", "    ")
		if err != nil {
			fmt.Println("Error converting to XML:", err)
			return
		}

		if err := sendDataToBank1(xmlData); err != nil {
			fmt.Println("Error sending data to server:", err)
			return
		}

	} else if id == "XML" {

	}

	// Receive updated data from server
	updatedDataXML, err := receiveDataFromServer()
	if err != nil {
		fmt.Println("Error receiving data from server:", err)
		return
	}

	// Convert received XML data to JSON
	var updatedData Data
	if err := xml.Unmarshal(updatedDataXML, &updatedData); err != nil {
		fmt.Println("Error parsing XML:", err)
		return
	}

	// Convert updated data to JSON
	updatedDataJSON, err := json.Marshal(updatedData)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}

	// Send updated data to API server
	if err := sendDataToAPIServer(updatedDataJSON); err != nil {
		fmt.Println("Error sending updated data to API server:", err)
		return
	}
}

func main2() {
	
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
	    fmt.Println("Error starting TCP server:", err)
	    return
	}
	defer ln.Close()

	fmt.Println("TCP server started. Listening on port 8080...")

	// Accept connections and handle them in a new goroutine
	for {
	    conn, err := ln.Accept()
	    if err != nil {
	        fmt.Println("Error accepting connection:", err)
	        break
	    }
	    fmt.Println("connected")
	    go handleConnection(conn)
	}
}

// Fetch data from C# API
func fetchDataFromAPI() ([]byte, error) {
	// Example: Fetching data from a C# API
	resp, err := http.Get("http://api.example.com/data")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)

}

// Forward data to another server
func sendDataToBank1(data []byte) error {
	// Example: Forwarding data to another server
	_, err := http.Post("http://destination-server.com/data", "application/xml", bytes.NewBuffer(data))
	return err
}

func receiveDataFromServer() ([]byte, error) {
	// Example: Receiving updated data from server
	resp, err := http.Get("http://server.example.com/updated-data")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Send updated data to API server
func sendDataToAPIServer1(data []byte) error {
	// Example: Sending updated data to API server
	_, err := http.Post("http://api.example.com/updated-data", "application/json", bytes.NewBuffer(data))
	return err
}
