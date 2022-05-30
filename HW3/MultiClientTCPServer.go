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
	"strings"
	"sync"
	"syscall"
	"time"
)

var numOfConnectedClients = 0
var numOfClients = 0
var countCommand = 0
var startTime time.Time
var mutex *sync.Mutex
var mutex2 *sync.Mutex
var wait sync.WaitGroup
var conList []net.Conn

func main() {
	startTime = time.Now()
	serverPort := "30768"
	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("Server is ready to receive on port %s\n", serverPort)
	var conn net.Conn
	var err error

	mutex = new(sync.Mutex)
	mutex2 = new(sync.Mutex)

	// check timer periodically
	go func() {
		timerTicker := time.NewTicker(60 * time.Second)
		for {
			select {
			case <-timerTicker.C:
				fmt.Printf("Number of connected clients = %d\n", numOfConnectedClients)
			}
		}
	}()

	// Disconnect and shut down on keyboard interrupt
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		fmt.Printf("\n")
		i := 0
		for i < len(conList) {
			if conList[i] != nil {
				conList[i].Close()
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
		fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())
		mutex.Lock()
		numOfClients += 1
		numOfConnectedClients += 1
		fmt.Printf("Client %d connected. Number of connected clients = %d\n", numOfClients, numOfConnectedClients)
		go processRequestClient(conn, numOfClients)
		wait.Add(1)
		mutex.Unlock()
	}
}

func checkConnectedClient(numC int, conn net.Conn) {
	mutex.Lock()
	numOfConnectedClients -= 1
	fmt.Printf("Client %d disconnected. Number of connected clients = %d\n", numC, numOfConnectedClients)
	wait.Done()
	conn.Close()
	mutex.Unlock()
}

func processRequestClient(conn net.Conn, numC int) {
	defer checkConnectedClient(numC, conn)
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
			countCommand += 1

		} else if string(buffer[0:1]) == "2" {
			point := strings.Index(conn.RemoteAddr().String(), ":")
			IP := conn.RemoteAddr().String()[:point]
			fmt.Printf("Command 2\n")
			port := conn.RemoteAddr().String()[point+1:]
			reply := "client IP = " + IP + ", port = " + port + "\n"
			conn.Write([]byte(reply))
			countCommand += 1

		} else if string(buffer[0:1]) == "3" {
			mutex2.Lock()
			conn.Write([]byte("requests served = " + fmt.Sprintf("%02d", countCommand) + "\n"))
			fmt.Printf("Command 3\n")
			countCommand += 1
			mutex2.Unlock()

		} else if string(buffer[0:1]) == "4" {
			elapsedTime := int(float64(time.Since(startTime)) / float64(time.Second))
			min := elapsedTime / 60
			hour := min / 60
			sec := elapsedTime % 60
			min = min % 60
			nowTime := fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)
			conn.Write([]byte("run time = " + nowTime + "\n"))
			fmt.Printf("Command 4\n")
			countCommand += 1
		}
	}

}
