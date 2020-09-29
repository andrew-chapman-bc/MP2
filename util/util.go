package util

import (
	
)

// Message holds username and message strings
type Message struct {
	Receiver string
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
	@function: CreateUserInputStruct
	@description: Uses a username and message to construct a Message struct
	@exported: True
	@params: string, string
	@returns: {Message}
*/
func CreateMessageStruct(receiver, message, sender string) Message {
	var input Message
	input.Receiver = receiver
	input.Message = message
	input.Sender = sender
	return input
}


// FormatMessage formats the message so it can be multi-spaced
func FormatMessage(messageArr []string) string {
	formattedMess := messageArr[2:]
	message := ""
	for _, word := range formattedMess {
		message += word
	}
	return message
}