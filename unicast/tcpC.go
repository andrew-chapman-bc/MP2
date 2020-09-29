package unicast

import (
	//"bufio"
	"fmt"
	//"log"
	"net"
	"os"
	"io/ioutil"
	"encoding/json"
	//"strings"
	//"time"
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




func ReadJSONForClient(userName string) Connections {
	jsonFile, err := os.Open("connections.json")
	var connections Connections
	if err != nil {
		fmt.Println(err)
		return connections
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	for i := 0; i < len(connections.Connections); i++ {
		if connections.Connections[i].Username == userName {
			return connections
		}
	}
	var newConn Connection
	newConn.Port = connections.Connections[0].Port
	newConn.Type = "client"
	newConn.Username = userName
	
	connections.Connections = append(connections.Connections, newConn)
	
	jsonData, err := json.Marshal(connections)
	if err != nil {
		fmt.Println("Error marshalling JSON")
		return connections
	}

	ioutil.WriteFile("connections.json", jsonData, os.ModePerm)
	return connections
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
func SendMessage(messageParams Message, connection Connections) {
	if messageParams.Message == "EXIT" {
		// TODO: close stuff
	} else {
		port := ""
		for _, connectionStruct := range connection.Connections {
			if connectionStruct.Type == "server" {
				port = connectionStruct.Port
			}
		}

		if port == "" {
			fmt.Println("Empty Port, message will not send.")
		}

		connectionString := connection.IP + ":" + port
		c, err := connectToTCPServer(connectionString)
		if (err != nil) {
			fmt.Println("Network Error: ", err)
		}
		jsonData, err := json.Marshal(messageParams)
		encoder := json.NewEncoder(c)
		_ = encoder.Encode(jsonData)
		

	}
	
	
} 

