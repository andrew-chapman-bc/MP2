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
)


// Server holds the structure of our TCP server implemenation
type Server struct {
    port   string
	server net.Listener
}

// NewTCPServer creates a new Server
func NewTCPServer(port string) (*Server, error) {
	server := Server{port: port}

	// if port is empty -> throw error
	if port == "" {
		return &server, errors.New("Error: Port not found")
	}

    return &server, nil
}


// RunServ starts the TCP Server.
func (serv *Server) RunServ(inputChan chan string, isUpdatingChan chan bool, wg *sync.WaitGroup) (err error) {
	// Create map of connections
	conns := make(map[string]net.Conn)

	serv.server, err = net.Listen("tcp", ":" + serv.port)
    if err != nil {
        return errors.New("Could not listen... Incorrect port?")
	}
    for {
        conn, err := serv.server.Accept()
        if err != nil || conn == nil {
            err = errors.New("Network Error: Can't accept this connection")
			break
		}
		serv.handleConnections(conns, isUpdatingChan, wg)
		// if we get an exit from I/O -> Close
		if inputData := <- inputChan; inputData == "EXIT" {
			conn.Close()
		}
    }
    return
}


func (serv *Server) handleConnections(conns map[string]net.Conn, isUpdatingChan chan bool, wg *sync.WaitGroup) (err error) {
	
    for {

		conn, err := serv.server.Accept()
		
        if err != nil || conn == nil {
            err = errors.New("Network Error: Could not accept connection")
            break
		}

		// constantly checking for new connections
		if isUpdating := <- isUpdatingChan; isUpdating == true {
			connections := serv.readJSONForServer(serv.port)
			for _, connect := range connections.Connections {
				if (connect.Type != "server") {
					conns[connect.Username] = conn
				}
			}
			isUpdatingChan <- false
		}
		wg.Add(1)
		defer wg.Done()
        go serv.handleConnection(conn, conns)
    }
    return
}

func (serv *Server) handleConnection(conn net.Conn, conns map[string]net.Conn) {

	dec := json.NewDecoder(conn)
	var mess util.Message
    for {
		err := dec.Decode(&mess)
		fmt.Println(mess)
		if err != nil {
			fmt.Println("Network Error: Error decoding JSON!")
			conn.Close()
		}

		serv.sendMessageToClient(mess, conns)
    }
}

func (serv *Server) sendMessageToClient(message util.Message, conns map[string]net.Conn) (err error) {
	
	conn := conns[message.Receiver]

	jsonData, err := json.Marshal(message)
	if err != nil {
		return errors.New("Error when marshalling json before client send")
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
	@function: readJSON
	@description: Reads the JSON and returns a struct which contains 
		the type, port, username and IP
	@exported: True
	@params: N/A
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

