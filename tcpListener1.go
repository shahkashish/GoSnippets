package main
import (

	"fmt"
	"log"
	"net"
	"bufio"
	"strings"
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"bytes"
)
type Data struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

type Server struct {
	listenAddress string
	ln            net.Listener
	quitch        chan struct{}
	msg           chan Message
}
type Message struct {
	from string
	payload []byte
}
func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddress: listenAddr,
        quitch: make(chan struct{}),
        msg: make(chan Message,10),
	}
}
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln
	go s.acceptLoop()
	
	<-s.quitch
	close(s.msg)
	return nil
}
func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error", err)
			continue
		}
		fmt.Println("new connection :", conn.RemoteAddr())
		go s.readLoop(conn)
	}
}
func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read error:", err)
			return
		}
		line = strings.TrimSpace(line)
		if line != "" {
			s.msg <- Message{
				from:    conn.RemoteAddr().String(),
				payload: []byte(line),
			}
		}
	}
}

func main() {
    server := NewServer(":8080")
    
    go func(){
        for msg := range(server.msg) {
            fmt.Printf("received message (%s): %s\n", msg.from, string(msg.payload))
            id, err := detectFormat(msg.payload)
            if err != nil {
                fmt.Println(err)
                continue
            }
            if id == "JSON" {
                var data Data
                if err := json.Unmarshal(msg.payload, &data); err != nil {
                    fmt.Println("Error parsing JSON:", err)
                    continue
                }

                // Convert to XML
                xmlData, err := xml.MarshalIndent(data, "", "    ")
                if err != nil {
                    fmt.Println("Error converting to XML:", err)
                    continue
                }

                if err := sendDataToBank(xmlData); err != nil {
                    fmt.Println("Error sending data to server:", err)
                    continue
                }

            } else if id == "XML" {
                var data Data
                if err := xml.Unmarshal(msg.payload, &data); err != nil {
                    fmt.Println("Error parsing XML:", err)
                    continue
                }

                // Convert to JSON
                jsonData, err := json.MarshalIndent(data, "", "    ")
                if err != nil {
                    fmt.Println("Error converting to JSON:", err)
                    continue
                }

                if err := sendDataToAPIServer(jsonData); err != nil {
                    fmt.Println("Error sending data to server:", err)
                    continue
                }
            }
        }
    }()
    log.Fatal(server.Start())
}


func detectFormat(data []byte) (string, error) {

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
func sendDataToAPIServer(data []byte) error {
	// Example: Sending updated data to API server
	_, err := http.Post("http://api.example.com/updated-data", "application/json", bytes.NewBuffer(data))
	return err
}
func sendDataToBank(data []byte) error {
	// Example: Forwarding data to another server
	_, err := http.Post("http://destination-server.com/data", "application/xml", bytes.NewBuffer(data))
	return err
}