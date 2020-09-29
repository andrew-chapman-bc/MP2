package tcp

import (
	"fmt"
	"errors"
	"net"
	"os"
	"io/ioutil"
	"encoding/json"
	"../util"
)





// Client holds the structure of our TCP Client implementation
type Client struct {
	Username string
	client net.Conn
}

// NewTCPClient creates a new Client
func NewTCPClient(username string) (*Client, error) {
	client := Client{Username: username}
	// if username is empty -> throw error
	if username == "" {
		return nil, errors.New("Error: Address not found")
	}

	return &client, nil
}

// RunCli starts the TCPClient
func (cli *Client) RunCli(messageChan chan util.Message) (err error) {
	connection, err := cli.readJSONForClient(cli.Username)
	if err != nil {
		return errors.New("Error when reading JSON on Client Side")
	}

	cli.client, err = net.Dial("tcp", connection.Port)
	if err != nil {
		fmt.Println(connection.Port)
		fmt.Println(err)
		return errors.New("Could not dial... Incorrect address?")
	}
	for {
		if messageData := <- messageChan; messageData.Message == "EXIT" {
			break
		}
		cli.sendMessageToServer(cli.client, messageChan)
	}
	return

}

func (cli *Client) sendMessageToServer(conn net.Conn, messageChan chan util.Message) (err error) {
	
	messageData := <- messageChan

	jsonData, err := json.Marshal(messageData)
	if err != nil {
		return errors.New("Error marshalling JSON Data")
	}

	encoder := json.NewEncoder(conn)
	encoder.Encode(jsonData)
	fmt.Println("data sent!", messageData)
	return
}


func (cli *Client) readJSONForClient(userName string) (util.Connection, error) {
	jsonFile, err := os.Open("connections.json")
	var connections util.Connections
	ourConnect := util.Connection{"","",""}
	if err != nil {
		return ourConnect, errors.New("Error opening JSON file on Client Side")
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	for i := 0; i < len(connections.Connections); i++ {
		if connections.Connections[i].Username == userName {
			connections.Connections[i].Port = connections.IP + ":" + connections.Connections[i].Port
			return connections.Connections[i], nil
		}
	}

	ourConnect.Port = connections.IP + ":" + connections.Connections[0].Port
	ourConnect.Type = "client"
	ourConnect.Username = userName
	
	connections.Connections = append(connections.Connections, ourConnect)
	
	jsonData, err := json.Marshal(connections)
	if err != nil {
		fmt.Println("Error marshalling JSON")
	}

	ioutil.WriteFile("connections.json", jsonData, os.ModePerm)
	return ourConnect, nil
}

