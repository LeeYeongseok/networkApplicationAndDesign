/**
 * 20170768 LeeYeongSeok
 **/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// truncate to 3 decimal places and return as a string
func cutFloat(num float64) string {
	return fmt.Sprintf("%06.3f", num)
}

// If there is an error in the server, notify and close and exit
func checkError(err error, conn net.Conn) {
	if nil != err {
		fmt.Println("server Disconnected")
		fmt.Println("Bye bye~")
		conn.Close()
		os.Exit(0)
	}
}

func main() {
	serverName := "nsl2.cau.ac.kr"
	serverPort := "30768"
	conn, err := net.Dial("tcp", serverName+":"+serverPort)
	if nil != err {
		fmt.Println("server Disconnected")
		fmt.Println("Bye bye~")
		os.Exit(0)
	}
	localAddr := conn.LocalAddr().(*net.TCPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

	// Disconnect and shut down on keyboard interrupt
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		fmt.Println("\nBye bye~")
		conn.Close()
		os.Exit(0)
	}()

	for {
		fmt.Printf("<Menu>\n1) convert text to UPPER-case\n2) get my IP address and port number\n3) get server request count\n4) get server running time\n5) exit\nInput option: ")
		num, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		buffer := make([]byte, 1024)
		var elapsedTime float64

		if num == "1\n" {
			fmt.Printf("Input sentence: ")
			input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			conn.Write([]byte("1" + input))
			//start
			startTime := time.Now()
			_, err := conn.Read(buffer)
			//finish
			elapsedTime = float64(time.Since(startTime)) / float64(time.Millisecond)
			checkError(err, conn)

		} else if num == "2\n" {
			conn.Write([]byte("2"))
			//start
			startTime := time.Now()
			_, err := conn.Read(buffer)
			//finish
			elapsedTime = float64(time.Since(startTime)) / float64(time.Millisecond)
			checkError(err, conn)

		} else if num == "3\n" {
			conn.Write([]byte("3"))
			//start
			startTime := time.Now()
			_, err := conn.Read(buffer)
			//finish
			elapsedTime = float64(time.Since(startTime)) / float64(time.Millisecond)
			checkError(err, conn)

		} else if num == "4\n" {
			conn.Write([]byte("4"))
			//start
			startTime := time.Now()
			_, err := conn.Read(buffer)
			//finish
			elapsedTime = float64(time.Since(startTime)) / float64(time.Millisecond)
			checkError(err, conn)

		} else if num == "5\n" {
			break
		} else {
			continue
		}

		fmt.Printf("\nReply from server: %s", string(buffer))
		fmt.Println("RTT = ", cutFloat(elapsedTime), "ms\n")

	}
	fmt.Println("Bye bye~")
	conn.Close()
}
