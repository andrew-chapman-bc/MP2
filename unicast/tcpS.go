package unicast

import (
	"fmt"
	"net"
	"bufio"
	"strings"
	"time"
	"os"
	"log"
	"io/ioutil"	
	"math/rand"
	"strconv"
	"encoding/json"
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
	@description: Uses a destination, message and source string to construct a UserInput struct to send and receive the message across the server/client
	@exported: True
	@params: string, string, string
	@returns: {UserInput}
*/
func CreateUserInputStruct(username, message string) UserInput {
	var input UserInput
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
func handleConnection(c net.Conn) {
	// even though we don't support multi-messaging at the moment, no reason to possibly be running this multiple times inside the for loop
	delay, err := getDelayParams()
	if (err != nil) {
		fmt.Println("Error: ", err)
	}
	
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
        if err != nil {
            fmt.Println(err)
            return
		}
		netArray := strings.Fields(netData)
		// generate the network delay on the receive side, must do it here and not in the sendmessage function because we are using goroutines
		generateDelay(delay)
		timeOfReceive := time.Now().Format("02 Jan 06 15:04:05.000 MST")
		fmt.Println("Received " + netArray[0] + " from destination " + netArray[1] + " system time is: " + timeOfReceive)
	}
}



/*
	@function: getDelayParams
	@description: Scans the config file for the first line to get the delay parameters that will be used to simulate the network delay
	@exported: false
	@params: N/A 
	@returns: Delay, error
*/
func getDelayParams() (Delay, error) {
	config, err := os.Open("config.txt")
	scanner := bufio.NewScanner(config)
	scanner.Split(bufio.ScanLines)
	success := scanner.Scan()
	if success == false {
		err = scanner.Err()
		if err == nil {
			fmt.Println("Scanned first line")
		} else {
			log.Fatal(err)
		}
	}
	delays := strings.Fields(scanner.Text())
	var delayStruct Delay
	delayStruct.minDelay = delays[0]
	delayStruct.maxDelay = delays[1]
	return delayStruct, err
} 

/*
	@function: generateDelay
	@description: Uses the delay parameters obtained from getDelayParams() to generate the delay that will be used in sendMessage function
	@exported: false
	@params: Delay
	@returns: N/A
*/
func generateDelay (delay Delay) {
	rand.Seed(time.Now().UnixNano())
	min, _ := strconv.Atoi(delay.minDelay)
	max, _ := strconv.Atoi(delay.maxDelay)
	delayTime := rand.Intn(max - min + 1) + min
	time.Sleep(time.Duration(delayTime) * time.Millisecond)
} 

/*
	@function: connectToTCPClient
	@description: Opens a concurrent TCP Server and calls the net.Listen function to connect to the TCP Client
	@exported: True
	@params: string
	@returns: N/A
*/
func ConnectToTCPClient(PORT string) {
	// listen/connect to the tcp client
	l, err := net.Listen("tcp4", ":" + PORT)
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConnection(c)
	}
}