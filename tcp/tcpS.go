package tcp

import (
	"fmt"
	"net"
	"os"
	"io/ioutil"	
	"encoding/json"
	"errors"
	"../util"
	"sync"
	"bufio"
)


// Server holds the structure of our TCP server implemenation
type Server struct {
    port   string
	server net.Listener
}


/*
	@function: NewTCPServer
	@description: Creates a Server Instance which can then be used in the main function
	@exported: True
	@family: N/A
	@params: string
	@returns: {*Server}, error
*/
func NewTCPServer(port string) (*Server, error) {
	server := Server{port: port}

	// if port is empty -> throw error
	if port == "" {
		return &server, errors.New("Error: Port not found")
	}

    return &server, nil
}

/*
	@function: RunServ
	@description: Starts the TCP server and calls handle connections
	@exported: True
	@family: Server
	@params: chan string, chan bool, waitgroup
	@returns: error
*/
func (serv *Server) RunServ(inputChan chan string, wg *sync.WaitGroup) (err error) {
	// Create map of connections
	conns := make(map[string]net.Conn)

	serv.server, err = net.Listen("tcp4", ":" + serv.port)
    if err != nil {
        return err
	}
	fmt.Println("Listening to the port:", serv.port)
	
	defer serv.server.Close()

    for {
		serv.handleConnections(conns, wg)
		inputData := <- inputChan
		if inputData == "EXIT" {
			serv.server.Close()
			break
		}
    }
    return
}
/*
	@function: handleConnections
	@description: calls the Accept function in a loop and calls another handleConnection goroutine which decodes data and sends it to the specified client
	@exported: false
	@family: Server
	@params: map[string]net.Conn, chan bool, WaitGroup
	@returns: error
*/
func (serv *Server) handleConnections(conns map[string]net.Conn, wg *sync.WaitGroup) (err error) {
	
	for {
		conn, err := serv.server.Accept()
		
        if err != nil || conn == nil {
            err = errors.New("Network Error: Could not accept connection")
            break
		}

		user, err := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("user", user)
		conns[user] = conn
		fmt.Println(conns)
		wg.Add(1)
        go serv.handleConnection(conn, conns)
	}
	wg.Done()
    return
}

/*
	@function: handleConnection
	@description: a goroutine which unserializes JSON data and then calls the sendMessageToClient function
	@exported: false
	@family: Server
	@params: net.Conn, map[string]net.Conn
	@returns: error
*/
func (serv *Server) handleConnection(conn net.Conn, conns map[string]net.Conn) (err error) {
	fmt.Println("ok ok")
	dec := json.NewDecoder(conn)
	mess := util.Message{"", "", ""}
    for {
		fmt.Println(mess)
		err := dec.Decode(&mess)
		if err != nil {
			fmt.Println(err)
			return err
		}
		mess := util.Message{}

		serv.sendMessageToClient(mess, conns)
    }
}

/*
	@function: sendMessageToClient
	@description: serializes JSON data and sends it over to the specified client
	@exported: false
	@family: Server
	@params: {Message}, map[string]net.Conn
	@returns: error
*/
func (serv *Server) sendMessageToClient(message util.Message, conns map[string]net.Conn) (err error) {
	
	conn := conns[message.Receiver]

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}
	if jsonData == nil {
		return errors.New("No JSON data before client send")
	}

	encoder := json.NewEncoder(conn)
	encoder.Encode(jsonData)
	fmt.Println("sent!")
	return
}



/*
	@function: readJSONForServer
	@description: Reads the JSON and returns a struct which contains 
		the type, port, username and IP
	@exported: False
	@family: Server
	@params: string
	@returns: Connections
*/
func (serv *Server) readJSONForServer(port string) util.Connections {
	jsonFile, err := os.Open("connections.json")
	var connections util.Connections
	if err != nil {
		fmt.Println(err)
		return connections
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &connections)
	for i := 0; i < len(connections.Connections); i++ {
		if (connections.Connections[i].Type == "server" ) {
			if (connections.Connections[i].Port != port) {
				connections.Connections[i].Port = port
			}
		}
	}
	jsonData, err := json.Marshal(connections)
	if err != nil {
		fmt.Println("Error marshalling JSON")
		return connections
	}
	// re-write to json
	ioutil.WriteFile("connections.json", jsonData, os.ModePerm)

	return connections
}

