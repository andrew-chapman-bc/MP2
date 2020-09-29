package main

import (
	"./unicast"
	"bufio"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
	"strings"
	"sync"
)

/*
	@function: getInput
	@description: gets the input entered through I/O and packages it into an array that will be used to create a {UserInput}
	@exported: False
	@params: N/A
	@returns: []string
*/
func getInput() []string {
	fmt.Println("Enter input >> ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	inputArray := strings.Fields(input)
	return inputArray

}

/*
	@function: parseInput
	@description: Parses the UserInput into a {Message}
	@exported: False
	@params: N/A
	@returns: {Message}
*/
func parseInput(message chan unicast.Message, sender string)  {
	inputArray := getInput()
	var inputStruct unicast.Message
	if inputArray[0] == "EXIT" {
		//messageString := "EXIT"
		inputStruct = unicast.CreateMessageStruct("", "EXIT", "")
		message <- inputStruct
	} else {
		messageString := unicast.FormatMessage(inputArray)
		inputStruct = unicast.CreateMessageStruct(inputArray[1], messageString, sender)
		message <- inputStruct
	}

}



/*
	@function: openTCPServerConnections
	@description:
	@exported: False
	@params: N/A
	@returns: N/A
*/
func openTCPServerConnections(connections unicast.Connections, message chan unicast.Message, wg *sync.WaitGroup, ) {
	serverConnection := unicast.Connection{}
	for i := 0; i < len(connections.Connections); i++ {
		if connections.Connections[i].Type == "server" {
			//used to have an err here, might need to put it back somewhere
			serverConnection = connections.Connections[i]
			break
		}
	}
	unicast.ConnectToTCPClient(serverConnection, message, wg)
}


func getCmdLine() string {
	// Use argparse library to get accurate command line data
	parser := argparse.NewParser("", "Concurrent TCP Channels")
	s := parser.String("s", "string", &argparse.Options{Required: true, Help: "String to print"})
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}
	return &s
}

func main() {
	
	s := getCmdLine()
	fmt.Println("WE GOT THE PARSE" , s)
	
	
	cmdLineArr := strings.Fields(s)
	connectionType := cmdLineArr[0]
	var wg sync.WaitGroup
	messageChan := make(chan unicast.Message)

	if (connectionType == "server") {
		port := cmdLineArr[1]
		connections := unicast.ReadJSONForServer(port)
		wg.Add(1)
		go openTCPServerConnections(connections, messageChan, &wg)
		messageData := <- messageChan
		wg.Add(1)
		go unicast.SendMessage(messageData, connections)
	} else {
		userName := cmdLineArr[1]
		connections := unicast.ReadJSONForClient(userName)
		parseInput(messageChan, userName)
		wg.Add(1)
		go func() {
			messageData := <- messageChan
			unicast.SendMessage(messageData, connections)
			wg.Done()
		}()
	}
	wg.Wait()
}

