/**
 * 20170768 LeeYeongSeok
 * TCPServer.go
 **/

package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

var nicknames []string
var mutex *sync.Mutex
var wait sync.WaitGroup
var conList []net.Conn
var version = "1.16.1"
var kick = false

func main() {
	serverPort := "30768"
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
		for _, data := range conList {
			tempList = append(tempList, data)
		}
		i := 0
		for i < len(tempList) {
			if tempList[i] != nil {
				tempList[i].Close()
			}
			i += 1
		}
		wait.Wait()
		listener.Close()
		fmt.Println("\nBye bye~")

		os.Exit(0)
	}()

	for {
		conn, err = listener.Accept()
		if nil != err {
			break
		}

		if len(nicknames) >= 3 { //over people
			conn.Write([]byte{0})
			conn.Close()
			continue
		} else { // not over people
			conn.Write([]byte{1})
		}

		go chattingRoom(conn)

		//fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())

		wait.Add(1)
	}
}

func checkConnectedClient(nickname string, conn net.Conn) {
	mutex.Lock()
	removeConn(nickname)
	disconnMessage := "<" + nickname + "> left. There are <" + strconv.Itoa(len(nicknames)) + "> users now"
	if kick {
		disconnMessage = "[" + nickname + " is disconnected. There are " + strconv.Itoa(len(nicknames)) + " users in the chat room.]"
		kick = false
	} else {
		sendMessageAllUser(disconnMessage, conn)
	}
	fmt.Println(disconnMessage)
	wait.Done()
	conn.Close()
	mutex.Unlock()
}

func chattingRoom(conn net.Conn) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if bytes.Compare(buffer[0:1], []byte{0}) != 0 {
		conn.Close()
		wait.Done()
		return
	}
	nickname := string(buffer[1:n])
	if contains(nicknames, nickname) { // duplicate nickName
		conn.Write([]byte{2})
		conn.Close()
		wait.Done()
		return
	} else { // not duplicate
		connect := []byte{3}
		connect = append(connect, []byte(strconv.Itoa(len(nicknames)+1))...)
		conn.Write(connect)
	}

	mutex.Lock()
	conList = append(conList, conn)
	nicknames = append(nicknames, nickname)
	mutex.Unlock()

	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
	fmt.Printf("%s joined from <%s:%d>. There are <%d> users connected\n", nickname, remoteAddr.IP, remoteAddr.Port, len(nicknames))

	defer checkConnectedClient(nickname, conn)
	if nil != err {
		return
	}

	for {
		n, err := conn.Read(buffer)
		if nil != err {
			break
		}

		if bytes.Compare(buffer[0:1], []byte{5}) == 0 { // message
			messagePart := string(buffer[1:n])
			messagePart = strings.ToLower(messagePart)
			prohibit := "i hate professor"
			message := nickname + "> " + string(buffer[1:n])
			if strings.Contains(messagePart, prohibit) {
				kickMessage := "[" + nickname + " is disconnected. There are " + strconv.Itoa(len(nicknames)-1) + " users in the chat room.]"
				kickBuffer := []byte{8}
				kickBuffer = append(kickBuffer, []byte(kickMessage)...)
				sendMessageToUser(kickBuffer, conn)
				sendMessageAllUser(message+"\n"+kickMessage, conn)
				kick = true
				conn.Close()
			} else {
				sendMessageAllUser(message, conn)
			}
		} else if bytes.Compare(buffer[0:1], []byte{1}) == 0 { // list
			userList := ""
			for index, thisCon := range conList {
				remoteAddr := thisCon.RemoteAddr().(*net.TCPAddr)
				userList += "<" + nicknames[index] + ", " + remoteAddr.IP.String() + ", " + strconv.Itoa(remoteAddr.Port) + ">\n"
			}
			message := []byte{4}
			message = append(message, []byte(userList)...)
			sendMessageToUser(message, conn)

		} else if bytes.Compare(buffer[0:1], []byte{2}) == 0 { // dm
			nameLen := int(buffer[1])
			receiveNickname := string(buffer[2 : 2+nameLen])
			dmMessage := string(buffer[2+nameLen : n])
			sendDmToNickName(dmMessage, receiveNickname, nickname)

		} else if bytes.Compare(buffer[0:1], []byte{3}) == 0 { // ver
			message := []byte{6}
			message = append(message, []byte(version)...)
			sendMessageToUser(message, conn)

		} else if bytes.Compare(buffer[0:1], []byte{4}) == 0 { // rtt
			message := []byte{7}
			sendMessageToUser(message, conn)
		}

	}
}

func removeConn(nickname string) {
	for n, element := range nicknames {
		if element == nickname {
			conList = append(conList[:n], conList[n+1:]...)
			nicknames = append(nicknames[:n], nicknames[n+1:]...)
		}
	}
}

func contains(slice []string, str string) bool {
	for _, element := range slice {
		if element == str {
			return true
		}
	}
	return false
}

func sendMessageToUser(message []byte, toConn net.Conn) {
	toConn.Write(message)
}

func sendDmToNickName(message string, nickName string, sendNickName string) {
	for n, name := range nicknames {
		if name == nickName {
			buffer := []byte{5}
			buffer = append(buffer, []byte("from "+sendNickName+": ")...)
			buffer = append(buffer, []byte(message)...)
			conList[n].Write(buffer)
		}
	}
}

func sendMessageAllUser(message string, myConn net.Conn) {
	for _, conn := range conList {
		if myConn == conn {
			continue
		}
		buffer := []byte{8}
		buffer = append(buffer, []byte(message)...)
		conn.Write(buffer)
	}
}

func processRequestClient(conn net.Conn, numC int) {
	conList = append(conList, conn)
	buffer := make([]byte, 1024)
	for {
		count, err := conn.Read(buffer)
		if nil != err {
			break
		}

		if string(buffer[0:1]) == "1" {
			conn.Write(bytes.ToUpper(buffer[1:count]))
			fmt.Printf("Command 1\n")
		}
	}

}
