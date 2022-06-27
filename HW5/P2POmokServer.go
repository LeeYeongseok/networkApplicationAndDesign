/**
 * 20170768 LeeYeongSeok
 * TCPServer.go
 **/

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var wait sync.WaitGroup
var connList []net.Conn
var UDPportList []string
var nickNameList []string

// var conn1 net.Conn
// var conn2 net.Conn
// var UDPport1 string
// var UDPport2 string
// var nickname1 string
// var nickname2 string
var mutex *sync.Mutex

func main() {
	serverPort := "50768"
	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)
	var conn net.Conn
	var err error

	mutex = new(sync.Mutex)
	// Disconnect and shut down on keyboard interrupt
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		fmt.Printf("\n")
		var tempList []net.Conn
		tempList = append(tempList, connList...)
		// tempList = append(tempList, conn1)
		// tempList = append(tempList, conn2)
		i := 0
		for i < len(tempList) {
			if tempList[i] != nil {
				tempList[i].Close()
			}
			i += 1
		}
		fmt.Println("\nBye bye~")
		wait.Wait()

		listener.Close()

		os.Exit(0)
	}()

	for {
		conn, err = listener.Accept()
		if nil != err {
			break
		}

		if len(nickNameList) == 0 { // 0 person
			wait.Add(1)
			go waitEmemy(conn)
		} else { // 1 person
			wait.Add(1)
			startGame(conn)

			conn1IP := connList[0].RemoteAddr().(*net.TCPAddr).IP.String()
			conn2IP := connList[1].RemoteAddr().(*net.TCPAddr).IP.String()
			messageData := make(map[string]interface{})
			messageData["UDPport"] = UDPportList[1]
			messageData["IP"] = conn2IP
			messageData["nickname"] = nickNameList[1]
			message, _ := json.Marshal(messageData)
			// message := append([]byte{1}, []byte(UDPportList[1])...)
			// message = append(message, []byte(conn2IP)...)
			fmt.Println(string(message))
			connList[0].Write(message)

			messageData2 := make(map[string]interface{})
			messageData2["UDPport"] = UDPportList[0]
			messageData2["IP"] = conn1IP
			messageData2["nickname"] = nickNameList[0]
			message2, _ := json.Marshal(messageData2)
			// message2 := append([]byte{1}, []byte(UDPportList[0])...)
			// message2 = append(message2, []byte(conn1IP)...)
			fmt.Println(string(message2))
			connList[1].Write(message2)

			clearUser()
		}
	}
}

func waitEmemy(conn net.Conn) {
	defer disconnectConn(conn)
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)
	nickname := string(buffer[6:n])
	UDPPort := string(buffer[1:6])
	connList = append(connList, conn)
	nickNameList = append(nickNameList, nickname)
	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
	fmt.Printf("%s joined from %s:%d. UDP port %s.\n", nickname, remoteAddr.IP, remoteAddr.Port, UDPPort)
	mutex.Lock()
	UDPportList = append(UDPportList, UDPPort)
	mutex.Unlock()
	fmt.Printf("1 user connected, waiting for another\n")
	conn.Write([]byte{0})

}

func startGame(conn net.Conn) {
	defer disconnectConn(conn)
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)
	nickname := string(buffer[6:n])
	UDPPort := string(buffer[1:6])
	connList = append(connList, conn)
	nickNameList = append(nickNameList, nickname)
	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
	mutex.Lock()
	UDPportList = append(UDPportList, UDPPort)
	mutex.Unlock()
	fmt.Printf("%s joined from %s:%d. UDP port %s.\n", nickname, remoteAddr.IP, remoteAddr.Port, UDPPort)
	fmt.Printf("2 user connected, notifying %s and %s\n", nickNameList[0], nickNameList[1])

}

func disconnectConn(conn net.Conn) {
	wait.Done()
}

func clearUser() {
	i := 0
	for i < len(connList) {
		if connList[i] != nil {
			connList[i].Close()
		}
		i += 1
	}
	nickNameList = nickNameList[:0]
	connList = connList[:0]
	UDPportList = UDPportList[:0]
}
