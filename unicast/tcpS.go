package unicast

import (
	"fmt"
	"net"
	"bufio"
	"strings"
	"os"
	"io/ioutil"	
	"encoding/json"
	"sync"
)



/*
	@function: readJSON
	@description: Reads the JSON and returns a struct which contains 
		the type, port, username and IP
	@exported: True
	@params: N/A
	@returns: Connections
*/
func ReadJSON() Connections {
	jsonFile, err := os.Open("connections.json")
	var connections Connections
	if err != nil {
		fmt.Println(err)
		return connections
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	return connections
}

/*
	@function: CreateUserInputStruct
	@description: Uses a username and message to construct a Message struct
	@exported: True
	@params: string, string
	@returns: {Message}
*/
func CreateMessageStruct(username, message string) Message {
	var input Message
	input.UserName = username
	input.Message = message
	return input
}


/*
	@function: handleConnection
	@description: handles connections to the concurrent TCP client and receives messages that are sent over through a goroutine in connectToTCPClient
	@exported: False
	@params: net.Conn
	@returns: N/A
*/
func handleConnection(c net.Conn, message chan Message) {

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
        if err != nil {
            fmt.Println(err)
            return
		}
		netArray := strings.Fields(netData)
		messageString := FormatMessage(netArray)
		// send username hey, im bored
		// STOP
		messageStruct := CreateMessageStruct(netArray[1], messageString)
		message <- messageStruct
	}
}

func FormatMessage(messageArr []string) string {
	formattedMess := messageArr[2:]
	message := ""
	for _, word := range formattedMess {
		message += word
	}
	return message
}


/*
	@function: connectToTCPClient
	@description: Opens a concurrent TCP Server and calls the net.Listen function to connect to the TCP Client
	@exported: True
	@params: string
	@returns: N/A
*/
func ConnectToTCPClient(connection Connection, message chan Message, wg *sync.WaitGroup) {
	port := connection.Port
	// listen/connect to the tcp client
	l, err := net.Listen("tcp4", ":" + port)
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()
	wg.Add(1)
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				fmt.Println(err)
			}
			handleConnection(c, message)
			messageStruct := <- message
			content := messageStruct.Message
			if content == "EXIT" {
				break
			}	
		}
		wg.Done()
	}()
}