package main


import (
	"fmt"
	"net"
	"bufio"
	"strings"
	"os"
	"io/ioutil"	
	"encoding/json"
	"github.com/akamensky/argparse"
	"sync"
	"encoding/gob"
	"strconv"
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



func getCmdLine() string {
	parser := argparse.NewParser("", "Private Chat Room")
	i := parser.Int("i", "int", &argparse.Options{Required:true, Help: "Source destination/identifiers"})
	err := parser.Parse(os.Args)
	s := strconv.Itoa(*i)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return s
	}
	return s
}

func disconnect(message chan Message) {
	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if input == "EXIT" {
			fmt.Println("TCP Server Exiting...")
			messageStruct := createMessageStruct("EXIT", "EXIT", "EXIT")
			message <- messageStruct
			return
		}
	}
}

/*
	@function: readJSON
	@description: Reads the JSON and returns a struct which contains 
		the type, port, username and IP
	@exported: True
	@params: N/A
	@returns: Connections
*/
func readJSON(port string) Connections {
	jsonFile, err := os.Open("connections.json")
	var connections Connections
	if err != nil {
		fmt.Println(err)
		return connections
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	for i := 0; i < len(connections.Connections); i++ {
		if (connections.Connections[i].Type == "server") {
			if connections.Connections[i].Port != port {
				connections.Connections[i].Port = port
				break
			}
		}
	}
	return connections
}

/*
	@function: CreateUserInputStruct
	@description: Uses a username and message to construct a Message struct
	@exported: True
	@params: string, string
	@returns: {Message}
*/
func createMessageStruct(username, message, sender string) Message {
	var input Message
	input.UserName = username
	input.Message = message
	input.Sender = sender
	return input
}


/*
	@function: handleConnection
	@description: handles connections to the concurrent TCP client and receives messages that are sent over through a goroutine in connectToTCPClient
	@exported: False
	@params: net.Conn
	@returns: N/A
*/
func handleConnection(c net.Conn) {
	fmt.Println("started 2")
	for {
		netData, _ := bufio.NewReader(c).ReadString('\n')
		netDataArr := strings.Fields(netData)
		var message Message
		message.UserName = netDataArr[0]
		message.Message = netDataArr[1] 
		message.Sender = netDataArr[2]

		decode := gob.NewDecoder(c)
		_ = decode.Decode(message)

		encode := gob.NewEncoder(c)
		encode.Encode(message)

	}
}

func formatMessage(messageArr []string) string {
	formattedMess := messageArr[2:]
	message := ""
	for _, word := range formattedMess {
		message += word
	}
	return message
}


/*
	@function: openTCPServerConnections
	@description:
	@exported: False
	@params: N/A
	@returns: N/A
*/
func openTCPServerConnections(connections Connections, message chan Message, wg *sync.WaitGroup) {
	var serverConnection Connection
	for i := 0; i < len(connections.Connections); i++ {
		if (connections.Connections[i].Type == "server") {
			serverConnection = connections.Connections[i]
			break
		}
	}
	fmt.Println(serverConnection)
	connectToTCPClient(serverConnection, message, wg)
}


/*
	@function: connectToTCPClient
	@description: Opens a concurrent TCP Server and calls the net.Listen function to connect to the TCP Client
	@exported: True
	@params: string
	@returns: N/A
*/
func connectToTCPClient(connection Connection, message chan Message, wg *sync.WaitGroup) {
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
			fmt.Println("hit")
			handleConnection(c)
			messageStruct := <- message
			content := messageStruct.Message
			if content == "EXIT" {
				break
			}	
		}
		wg.Done()
	}()
}



func main() {
	s := getCmdLine()
	var wg sync.WaitGroup
	connections := readJSON(s)
	message := make(chan Message)
	wg.Add(2)
	go openTCPServerConnections(connections, message, &wg)
	go disconnect(message)
	fmt.Println("Started")
	wg.Wait()

}
