/**
 * 20170768 LeeYeongSeok
 **/

package main

import (
	"bufio"
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

var startTime time.Time
var mutex2 *sync.Mutex

// truncate to 3 decimal places and return as a string
func cutFloat(num float64) string {
	return fmt.Sprintf("%06.3f", num)
}

// If there is an error in the server, notify and close and exit
func checkError(err error, conn net.Conn) {
	if nil != err {
		fmt.Println("gg~")
		conn.Close()
		os.Exit(0)
	}
}

func IsLetter(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}

func main() {
	serverName := "nsl2.cau.ac.kr"
	serverPort := "30768"
	if len(os.Args) < 2 {
		fmt.Println("please input nickname!")
		return
	}
	if len(os.Args) > 2 {
		fmt.Println("don't use blank in nickname!")
		return
	}
	nickName := os.Args[1]

	// Disconnect and shut down on keyboard interrupt
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		fmt.Println("\ngg~")
		os.Exit(0)
	}()

	if len(nickName) > 32 {
		fmt.Println("too long nickName")
		fmt.Println("Bye bye~")
		os.Exit(0)
	}
	if !IsLetter(nickName) {
		fmt.Println("You can use nickname only english")
		fmt.Println("Bye bye~")
		os.Exit(0)
	}
	conn, err := net.Dial("tcp", serverName+":"+serverPort)
	if nil != err {
		fmt.Println("server Disconnected")
		fmt.Println("gg~")
		os.Exit(0)
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if bytes.Compare(buffer[0:1], []byte{0}) == 0 {
		fmt.Printf("[chatting room full. cannot connect]\n")
		conn.Close()
		os.Exit(0)
	}
	pck := append([]byte{0}, []byte(nickName)...)
	_, err = conn.Write(pck)
	checkError(err, conn)

	n, err = conn.Read(buffer)
	if bytes.Compare(buffer[0:1], []byte{2}) == 0 {
		fmt.Printf("[that nickname is already used by another user. cannot connect.]\n")
		conn.Close()
		os.Exit(0)
	} else if bytes.Compare(buffer[0:1], []byte{3}) == 0 {
		remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
		fmt.Printf("[welcome %s to CAU network class chat room at %s:%d.] \n[There are %s users connected.]\n", nickName, remoteAddr.IP, remoteAddr.Port, buffer[1:n])
	}

	go func() {
		<-signals
		conn.Close()
	}()

	go listenTCP(conn)

	for {
		sc := bufio.NewScanner(os.Stdin)
		sc.Scan()
		message := sc.Text()
		if len(message) == 0 {
			continue
		}

		if message[0:1] == "\\" {
			if message == "\\exit" {
				break
			} else if message == "\\list" {
				_, err = conn.Write([]byte{1})
				checkError(err, conn)
			} else if len(message) > 5 && message[0:4] == "\\dm " {
				slice := strings.Split(message, " ")
				buffer := []byte{2}
				buffer = append(buffer, byte(len(slice[1])))
				buffer = append(buffer, []byte(slice[1])...)
				buffer = append(buffer, []byte(message[5+len(slice[1]):])...)
				_, err = conn.Write(buffer)
				checkError(err, conn)
			} else if message == "\\ver" {
				_, err = conn.Write([]byte{3})
				checkError(err, conn)
			} else if message == "\\rtt" {
				startTime = time.Now()
				_, err = conn.Write([]byte{4})
				checkError(err, conn)
			} else {
				fmt.Println("[invalid command]")
				continue
			}

		} else {
			pck := []byte{5}
			pck = append(pck, message...)
			_, err = conn.Write(pck)
			checkError(err, conn)
		}
	}
	conn.Close()
}

func listenTCP(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		checkError(err, conn)

		if bytes.Compare(buffer[0:1], []byte{8}) == 0 { // message
			fmt.Println("\n" + string(buffer[1:n]) + "\n")
		} else if bytes.Compare(buffer[0:1], []byte{4}) == 0 { // get list
			fmt.Println(string(buffer[1:n]))
		} else if bytes.Compare(buffer[0:1], []byte{5}) == 0 { // get dm
			fmt.Println(string(buffer[1:n]) + "\n")
		} else if bytes.Compare(buffer[0:1], []byte{6}) == 0 { // get version
			fmt.Println(string(buffer[1:n]) + "\n")
		} else if bytes.Compare(buffer[0:1], []byte{7}) == 0 { // get rtt
			elapsedTime := float64(time.Since(startTime)) / float64(time.Millisecond)
			fmt.Println("[RTT = ", cutFloat(elapsedTime), "ms]\n")
		}

	}
}
