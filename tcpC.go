package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
	"encoding/json"
	"io/ioutil"
)

// Message holds username and message strings
type Message struct {
	UserName string
	Message  string
	Sender   string
}

/*
	Connections: []Connection
	IP: IP Address to connect to
*/
type Connections struct {
	Connections []Connection `json:"connections"`
	IP string `json:"IP"`
} 

/*
	Type: "Server"/"Client" whether it's server or client
	Port: "1234", etc. Port attached to username
	Username: name of connection
	IP: IP address to connect to
*/
type Connection struct {
	Type string `json:"Type"`
	Port string `json:"Port"`
	Username string `json:"Username"`
}

/*
	@function: connectToTCPServer
	@description:	Connects to the TCP server with the ip/port obtained from config file as a parameter and 
					returns the connection to the server which will later be used to write to the server
	@exported: false
	@params: string 
	@returns: net.Conn, err
*/
func connectToTCPServer(connect string) (net.Conn, error) {
	// Dial in to the TCP Server, return the connection to it
	c, err := net.Dial("tcp", connect)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	return c, err
} 

/*
	@function: SendMessage
	@description: 	SendMessage sends the message from TCPClient to TCPServer by connecting to the server and 
					using the Fprintf function to send the message.
	@exported: True
	@params: {UserInput}, {Connection}
	@returns: N/A
*/
func sendMessage(messageParams chan Message, connection Connections, s string) {
	message := <- messageParams
	if (message.Message == "EXIT") {
		
	} else {
		port := ""
		for _, connectionStruct := range connection.Connections {
			if (connectionStruct.Type == "server") {
				port = connectionStruct.Port
			}
		}

		if port == "" {
			fmt.Println("Empty Port, message will not send.")
		}
		ip := connection.IP
		connectionString := ip + ":" + port
		c, err := connectToTCPServer(connectionString)
		if (err != nil) {
			fmt.Println("Network Error: ", err)
		}
		fmt.Fprintf(c, message.UserName + " " + message.Message + " " + s)

	}
	
} 

func readJSONForClient(userName string) Connections {
	jsonFile, err := os.Open("connections.json")
	var connections Connections
	if err != nil {
		fmt.Println(err)
		return connections
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	serverPort := ""
	for i := 0; i < len(connections.Connections); i++ {
		if (connections.Connections[i].Type == "server") {
			serverPort = connections.Connections[i].Port
		}
		if (connections.Connections[i].Username == userName) {
			break
		}
	}
	var connection Connection
	connection.Port = serverPort
	connection.Type = "client"
	connection.Username = userName
	connections.Connections = append(connections.Connections, connection)
	newJSON, _ := json.Marshal(connections)
	ioutil.WriteFile("connections.json", newJSON, os.ModePerm)
	return connections

}
/*
	@function: parseInput
	@description: Parses the UserInput into a {Message}
	@exported: False
	@params: N/A
	@returns: {Message}
*/
func parseInput(message chan Message, s string)  {
	inputArray := getInput()
	var inputStruct Message
	if inputArray[0] == "EXIT" {
		messageString := "EXIT"
		inputStruct = createMessageStruct("", "EXIT", "")
		message <- inputStruct
	} else {
		messageString := formatMessage(inputArray)
		inputStruct = createMessageStruct(inputArray[1], messageString, s)
		message <- inputStruct
	}

}
// send andy this is a message

func getInput() []string {
	fmt.Println("Enter input >> ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	inputArray := strings.Fields(input)
	return inputArray

}

func main() {
	s := getCmdLine()
	if *s == "" {
		fmt.Println("Error: Enter valid Username")
	}
	messageChan := make(chan Message)

	connections := readJSONForClient(*s)

	for {
		parseInput(messageChan, *s)
		go sendMessage(messageChan, connections, *s)
		message := <- messageChan
		if (message.Message == "EXIT") {
			break
		}
	}
	
}
