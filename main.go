package main

import (
	"./tcp"
	"bufio"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
	"strings"
	"sync"
	"errors"
	"./util"
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
func parseServInput(inputChan chan string) (err error) {
	inputArray := getInput()
	if len(inputArray) != 0 {
		if inputArray[0] == "EXIT" {
			inputChan <- "EXIT"
		}
	} else {
		return errors.New("Error parsing server input")
	}
	return
}

func parseCliInput(messageChan chan util.Message, userName string) (err error) {
	messageArray := getInput()
	if len(messageArray) != 0 {
		messageStruct := util.CreateMessageStruct(messageArray[1], util.FormatMessage(messageArray), userName)
		messageChan <- messageStruct
	} else {
		return errors.New("Error parsing client input")
	}
	return
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
	return *s
}

func main() {
	
	s := getCmdLine()
	cmdLineArr := strings.Fields(s)
	connectionType := cmdLineArr[0]
	var wg sync.WaitGroup
	inputChan := make(chan string)
	isUpdatingChan := make(chan bool)
	messageChan := make(chan util.Message)
	
	// Server -> Read JSON/Read Input/handleConnections 
	// Client -> Write JSON/Read Input/Write to server
	switch connectionType {
	case "server":
		var serv *tcp.Server
		port := cmdLineArr[1]
		serv, err := tcp.NewTCPServer(port)
		if err != nil {
			fmt.Println(err)
		}
		wg.Add(2)
		go func() {
			defer wg.Done()
			err := serv.RunServ(inputChan, isUpdatingChan, &wg)
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
		go func() {
			defer wg.Done()
			err := parseServInput(inputChan)
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
	default:
		var cli *tcp.Client
		userName := cmdLineArr[1] 
		cli, err := tcp.NewTCPClient(userName)
		if err != nil {
			fmt.Println(err)
			break
		}
		wg.Add(2)
		go func() {
			defer wg.Done()
			err := cli.RunCli(messageChan)
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
		go func() {
			defer wg.Done()
			err := parseCliInput(messageChan, userName)
			if err != nil {
				fmt.Println(err)
				return
			}
		}()

	} 
	wg.Wait()
}

