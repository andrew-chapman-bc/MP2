package unicast

import (
	//"bufio"
	"fmt"
	//"log"
	"net"
	//"os"
	//"strings"
	//"time"
)

// Message holds username and message strings
type Message struct {
	UserName string
	Message  string
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
func SendMessage(messageParams Message, connection Connections, ip string) {
	if messageParams.Message == "EXIT" {
		
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

		connectionString := ip + ":" + port
		c, err := connectToTCPServer(connectionString)
		if (err != nil) {
			fmt.Println("Network Error: ", err)
		}
		fmt.Fprintf(c, messageParams.Message)

	}
	
	// Sending the message to TCP Server
	// Easier to send this over as strings since it is only one message, we want the source to know where it comes from
	//fmt.Fprintf(c, messageParams.Message + " " + "\n")
	//timeOfSend := time.Now().Format("02 Jan 06 15:04:05.000 MST")
	//fmt.Println("Sent message " + messageParams.Message + " to destination " + messageParams.UserName + " system time is: " + timeOfSend)
	
} 

