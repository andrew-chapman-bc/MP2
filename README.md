# MP2
PM Chat Room
--- 
# To Run

In terms of dependencies, the only one that you need to download is the argparse golang package.  You can download this by running this command into one of your open terminals:
```
$ go get -u -v github.com/akamensky/argparse
``` 

Simply repeat steps in new terminals to have as many processes as needed

Open one terminal and enter
```bash
go run main.go --string "server 1234 server"
``` 
Open up a second terminal and enter
```bash
go run main.go --string "client 1234 andrew"
``` 
Then open up a third terminal and enter
```bash
go run main.go --string "client 1234 matt"
```
To send a message, go to second terminal and send

```bash
send matt hello there
```

Should output on third terminal
```bash 
Received hello there from andrew
```

If you want to send a message back to terminal 2, input
```bash
send andrew hi sir
```

The  following output should be printed on terminal 2 
```bash
Received hi sir from matt
```

---

# Structure and design

TCP Server
The first terminal made is a concurrent TCP server
The user's commandline input decides whether it is a server or client, as well as port number

The commandline input is written into the connections.json file
We have the ip hardcoded since this is all local host, but it can be added to the connections field
in a more complex scenario 


We have three structures designed to make passing data around easier and more readable.

They are as follows: 

The Message struct is used to easily access and pass around the username and message (send what to where) throughout the program without sending strings everywhere

```
type Message struct {
	UserName string
	Message  string
}

```

The Connections struct has the ip an array of username, port, client/server status of all terminals accessable in our codebase
(An array of connection structs) 
```
type Connections struct {
	Connections []Connection `json:"connections"`
	IP string `json:"IP"`
} 
```

Stores whether terminal is client or server, port, and username
```
type Connection struct {
	Type string `json:"Type"`
	Port string `json:"Port"`
	Username string `json:"Username"`
}
```
The ladder two structs use the json file to get their information 

# json file
The json file has the following format 
-----------------------------------------------------------------------------------------------
```    
"connections": [
        {
            "Type": "xxx",
            "Port": "xxx",
            "Username": "xxx"
        },
        {
            "Type": "xxx",
            "Port": "xxx",
            "Username": "xxx"
        }
    ],
    "IP": "127.0.0.1"
```
.... .... .......
-----------------------------------------------------------------------------------------------
To read the json file, there are two functions.
One function for the server reading, and one for the client


To add more connection, simply open a new terminal and run the program

For example:
-----------------------------------------------------------------------------------------------  
```  
{
    "connections": [
        {
            "Type": "server",
            "Port": "1234",
            "Username": "Matt"
        },
        {
            "Type": "client",
            "Port": "4567",
            "Username": "Andy"
        },
        {
            "Type": "client",
            "Port": "8543",
            "Username": "Lewis"
        }
    ],
    "IP": "127.0.0.1"
}


```
-----------------------------------------------------------------------------------------------

Goes to 

-----------------------------------------------------------------------------------------------   
``` 
{
    "connections": [
        {
            "Type": "server",
            "Port": "1234",
            "Username": "Matt"
        },
        {
            "Type": "client",
            "Port": "4567",
            "Username": "Andy"
        },
        {
            "Type": "client",
            "Port": "8543",
            "Username": "Lewis"
        },
        {
            "Type": "client",
            "Port": "1432",
            "Username": "Darius"
        }
    ],
    "IP": "127.0.0.1"
}

```
-----------------------------------------------------------------------------------------------

If you run the program again with the username "Darius"


# Input
The user input is broken up into three sections, : 
1. "Send"
2. Username 
3. Message

The program reads each section as follows: 
1. Disregard this keyword
2. Store the username into Message struct 
3. Store the message into Message struct

We are communicating between the server and client using channels 

# Exit Condition 

If the user enters "EXIT" the program will terminate its connection
The user will no longer be able to send/recieve messages

# Processes
The processes can be found in the unicast directory

In tcpC.go is where the message is sent out to the server

In tcpS.go is where the message is handled


### Shortcomings and Potential Improvemnts 
As of right now, the program does not run

