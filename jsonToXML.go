package main

import (
    "encoding/json"
    "encoding/xml"
    "fmt"
)

type Data2 struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
    Email string `json:"email"`
}

func main1() {
    // Sample JSON data received from the C# API
    jsonData := `{"name": "John", "age": 30, "email": "john@example.com"}`

    // Parse JSON data
    var data Data
    if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
        fmt.Println("Error parsing JSON:", err)
        return
    }

    // Convert to XML
    xmlData, err := xml.MarshalIndent(data, "", "    ")
    if err != nil {
        fmt.Println("Error converting to XML:", err)
        return
    }

    fmt.Println(string(xmlData))
}
