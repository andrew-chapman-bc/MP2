package main
import (
	"./unicast"
	"bufio"
	"fmt"
	"os"
	"sync"
	"strings"

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
	@description: Parses the UserInput into a {UserInput} and calls ScanConfig() to parse the parameters of TCP connection into a {Connection}
	@exported: False
	@params: N/A
	@returns: {UserInput}, {Connection}
*/
func parseInput(source *string) (unicast.UserInput, unicast.Connection) {
	inputArray := getInput()
	inputStruct := unicast.CreateUserInputStruct(inputArray[1], inputArray[2], *source)
	connection := unicast.ScanConfigForClient(inputStruct)
	return inputStruct, connection
}


// TODO: Scan for json here
/*
	@function: openTCPServerConnections
	@description: Opens all of the ports defined in the config file using ScanConfigForServer() to get an array of ports 
					and ConnectToTCPClient() to open them
	@exported: False
	@params: {WaitGroup}
	@returns: N/A
*/
func openTCPServerConnections(source *string) {
	// Need to send the source string in here so we know what port to look for
	openPort, err := unicast.ScanConfigForServer(*source)
	if openPort == "" {
		fmt.Println(err)
	}
	unicast.ConnectToTCPClient(openPort)
}

/*
	@function: unicast_send
	@description: function used as a goroutine to call SendMessage() to pass data from client to server, utilizes waitgroup
	@exported: False
	@params: {UserInput}, {Connection}, {WaitGroup}
	@returns: N/A
*/
func unicastSend(inputStruct unicast.UserInput, connection unicast.Connection) {
	// Send the message using UserInput struct and Connection struct to easily pass around data
	unicast.SendMessage(inputStruct, connection)
}

func main() {
	
	connections := unicast.ReadJSON()
	fmt.Println(connections)

	// // Use a wait group for goroutines
	// var wg sync.WaitGroup
	// // TODO: Scan for json => anonymous functions with wg.Done() + incorporate channels
	// wg.Add(1)
	// go openTCPServerConnections(&s)
	// inputStruct, connection := parseInput(&s)
	// wg.Add(1)
	// go unicastSend(inputStruct, connection, &wg)
	// wg.Wait()
}