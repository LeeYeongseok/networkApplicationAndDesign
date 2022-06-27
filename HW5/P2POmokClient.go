/**
 * 20170768 LeeYeongSeok
 **/

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var opUDPport string
var opIP string
var opNickName string
var nickName string
var myturn int
var gameFinish int

type Board [][]int

const (
	Row = 10
	Col = 10
)

func printBoard(b Board) {
	fmt.Print("   ")
	for j := 0; j < Col; j++ {
		fmt.Printf("%2d", j)
	}

	fmt.Println()
	fmt.Print("  ")
	for j := 0; j < 2*Col+3; j++ {
		fmt.Print("-")
	}

	fmt.Println()

	for i := 0; i < Row; i++ {
		fmt.Printf("%d |", i)
		for j := 0; j < Col; j++ {
			c := b[i][j]
			if c == 0 {
				fmt.Print(" +")
			} else if c == 1 {
				fmt.Print(" 0")
			} else if c == 2 {
				fmt.Print(" @")
			} else {
				fmt.Print(" |")
			}
		}

		fmt.Println(" |")
	}

	fmt.Print("  ")
	for j := 0; j < 2*Col+3; j++ {
		fmt.Print("-")
	}

	fmt.Println()
}

func checkWin(b Board, x, y int) int {
	lastStone := b[x][y]
	startX, startY, endX, endY := x, y, x, y

	// Check X
	for startX-1 >= 0 && b[startX-1][y] == lastStone {
		startX--
	}
	for endX+1 < Row && b[endX+1][y] == lastStone {
		endX++
	}

	if endX-startX+1 >= 5 {
		return lastStone
	}

	// Check Y
	startX, startY, endX, endY = x, y, x, y
	for startY-1 >= 0 && b[x][startY-1] == lastStone {
		startY--
	}
	for endY+1 < Row && b[x][endY+1] == lastStone {
		endY++
	}

	if endY-startY+1 >= 5 {
		return lastStone
	}

	// Check Diag 1
	startX, startY, endX, endY = x, y, x, y
	for startX-1 >= 0 && startY-1 >= 0 && b[startX-1][startY-1] == lastStone {
		startX--
		startY--
	}
	for endX+1 < Row && endY+1 < Col && b[endX+1][endY+1] == lastStone {
		endX++
		endY++
	}

	if endY-startY+1 >= 5 {
		return lastStone
	}

	// Check Diag 2
	startX, startY, endX, endY = x, y, x, y
	for startX-1 >= 0 && endY+1 < Col && b[startX-1][endY+1] == lastStone {
		startX--
		endY++
	}
	for endX+1 < Row && startY-1 >= 0 && b[endX+1][startY-1] == lastStone {
		endX++
		startY--
	}

	if endY-startY+1 >= 5 {
		return lastStone
	}

	return 0
}

func clear() {
	fmt.Printf("%s", runtime.GOOS)

	clearMap := make(map[string]func()) //Initialize it
	clearMap["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clearMap["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	value, ok := clearMap[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                             //if we defined a clearMap func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clearMap terminal screen :(")
	}
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
	//serverName := "nsl2.cau.ac.kr"
	serverName := "192.168.0.7"
	serverPort := "50768"

	pconn, _ := net.ListenPacket("udp", ":")
	localAddr := pconn.LocalAddr().(*net.UDPAddr)
	if len(os.Args) < 2 {
		fmt.Println("please input nickname!")
		return
	}
	if len(os.Args) > 2 {
		fmt.Println("don't use blank in nickname!")
		return
	}
	nickName = os.Args[1]

	// Disconnect and shut down on keyboard interrupt
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		if pconn != nil {
			if gameFinish == 0 {
				fmt.Println("you lose!")
			}
			message := []byte{3}
			op_addr, _ := net.ResolveUDPAddr("udp", opIP+":"+opUDPport)
			pconn.WriteTo(message, op_addr)
			pconn.Close()
		}
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

	//disconnect TCP keyboard Interrupt
	go func() {
		<-signals
		conn.Close()
	}()

	// TCP part
	buffer := make([]byte, 1024)
	myUDPPort := localAddr.Port
	myUDPPortStr := strconv.Itoa(myUDPPort)
	if len(myUDPPortStr) < 5 {
		myUDPPortStr = "0" + myUDPPortStr
	}
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(myUDPPort))

	pck := append([]byte{0}, []byte(myUDPPortStr)...)
	pck = append(pck, []byte(nickName)...)

	_, err = conn.Write(pck)
	checkError(err, conn)

	n, err := conn.Read(buffer)
	checkError(err, conn)
	if bytes.Equal(buffer[0:1], []byte{0}) { // nobody wait at server
		remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
		fmt.Printf("welcome %s to p2p-omok server at %s:%d. \nwaiting for an opponent\n", nickName, remoteAddr.IP, remoteAddr.Port)

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		checkError(err, conn)
		var data map[string]interface{}
		json.Unmarshal(buffer[:n], &data)
		opUDPport = data["UDPport"].(string)
		opIP = data["IP"].(string)
		opNickName = data["nickname"].(string)

		myturn = 0

	} else {
		remoteAddr := conn.RemoteAddr().(*net.TCPAddr)
		var data map[string]interface{}
		json.Unmarshal(buffer[:n], &data)
		opUDPport, _ = data["UDPport"].(string)
		opIP = data["IP"].(string)
		opNickName = data["nickname"].(string)

		fmt.Printf("welcome %s to p2p-omok server at %s:%d. \n%s is waiting for you\n", nickName, remoteAddr.IP, remoteAddr.Port, opNickName)
		myturn = 1
	}
	conn.Close()

	// UDP part
	op_addr, _ := net.ResolveUDPAddr("udp", opIP+":"+opUDPport)

	fmt.Println("game start")
	time.Sleep(1 * time.Second)

	gameFinish = 0

	board := Board{}
	x, y, turn, count, win := -1, -1, 0, 0, 0
	for i := 0; i < Row; i++ {
		var tempRow []int
		for j := 0; j < Col; j++ {
			tempRow = append(tempRow, 0)
		}
		board = append(board, tempRow)
	}
	clear()
	printBoard(board)

	played := make(chan bool)

	// listen part
	go func() {
		for {
			buffer := make([]byte, 1024)
			n, _, _ := pconn.ReadFrom(buffer)

			if bytes.Equal(buffer[0:1], []byte{4}) { // chatting
				fmt.Printf("%s> %s\n", opNickName, string(buffer[1:n]))

			} else if bytes.Equal(buffer[0:1], []byte{1}) { // play omok
				x1, _ := strconv.Atoi(string(buffer[1]))
				y1, _ := strconv.Atoi(string(buffer[2]))
				if turn == 0 {
					board[x1][y1] = 1
				} else {
					board[x1][y1] = 2
				}

				clear()
				printBoard(board)

				win = checkWin(board, x1, y1)
				if win != 0 {
					if win == myturn+1 {
						fmt.Printf("you win!\n")
					} else {
						fmt.Printf("you lose!\n")
					}

					gameFinish = 1
					continue

				}

				count += 1
				if count == Row*Col {
					fmt.Printf("draw!\n")
					gameFinish = 1
					continue
					//break
				}

				turn = (turn + 1) % 2

			} else if bytes.Equal(buffer[0:1], []byte{2}) { // gg
				fmt.Printf("you win!\n")
				gameFinish = 1
			} else if bytes.Equal(buffer[0:1], []byte{3}) { // exit
				fmt.Printf("%s out!\n", opNickName)
				if gameFinish == 0 {
					fmt.Printf("you win!\n")
				}
				fmt.Printf("gg~\n")
				pconn.Close()
				os.Exit(0)
			}
		}
	}()

	// time count
	go func() {
		timeOutCount := 0
		for {
			if gameFinish == 1 {
				break
			}
			if turn == myturn {
				timeoutChan := time.After(10 * time.Second)
				select {
				case <-timeoutChan:
					if gameFinish == 1 {
						break
					}
					fmt.Printf("timeOut\n")
					fmt.Printf("you lose!\n")
					gameFinish = 1
					message := []byte{2}
					pconn.WriteTo(message, op_addr)
					break
				case <-played:
					timeOutCount += 1
				}
			}
		}
	}()

	// write part
	for {

		sc := bufio.NewScanner(os.Stdin)
		sc.Scan()
		input := sc.Text()
		if len(input) == 0 {
			continue
		} else if input[0:1] == "\\" { // command
			if input == "\\gg" {
				if gameFinish == 0 {
					message := []byte{2}
					pconn.WriteTo(message, op_addr)
					fmt.Printf("you lose!\n")
					gameFinish = 1
				} else {
					fmt.Printf("game finished!\n")
				}
				continue

			} else if input == "\\exit" {
				message := []byte{3}
				pconn.WriteTo(message, op_addr)
				if gameFinish == 0 {
					fmt.Printf("you lose!\n")
				}
				fmt.Printf("Bye~\n")
				break
			} else if len(input) <= 2 {
				fmt.Println("invalid command")
				continue
			} else if input[0:2] == "\\\\" {
				if gameFinish == 1 {
					fmt.Println("game finished")
					continue
				}
				if turn != myturn {
					fmt.Println("not your turn")
					continue
				}
				slice := strings.Split(input, " ")
				if len(slice) != 3 {
					fmt.Println("error, must enter x y!")
					time.Sleep(1 * time.Second)
					continue
				}
				if _, err := strconv.Atoi(slice[1]); err != nil {
					fmt.Printf("%q doesn't look like number.\n", slice[1])
					continue
				}
				if _, err := strconv.Atoi(slice[2]); err != nil {
					fmt.Printf("%q doesn't look like number.\n", slice[2])
					continue
				}
				x, _ = strconv.Atoi(slice[1])
				y, _ = strconv.Atoi(slice[2])
				if x < 0 || y < 0 || x >= Row || y >= Col {
					fmt.Println("error, out of bound!")
					time.Sleep(1 * time.Second)
					continue
				} else if board[x][y] != 0 {
					fmt.Println("error, already used!")
					time.Sleep(1 * time.Second)
					continue
				}

				message := []byte{1}
				message = append(message, []byte(strconv.Itoa(x))...)
				message = append(message, []byte(strconv.Itoa(y))...)
				pconn.WriteTo(message, op_addr)

			} else {
				fmt.Println("invalid command")
				continue
			}
		} else { // chatting
			message := []byte{4}
			message = append(message, []byte(input)...)
			pconn.WriteTo(message, op_addr)
			continue
		}

		if turn == 0 {
			board[x][y] = 1
		} else {
			board[x][y] = 2
		}

		clear()
		printBoard(board)

		win = checkWin(board, x, y)
		if win != 0 {
			if win == myturn+1 {
				fmt.Printf("you win!\n")
			} else {
				fmt.Printf("you lose!\n")
			}
			gameFinish = 1
			continue
			//break
		}

		count += 1
		if count == Row*Col {
			fmt.Printf("draw!\n")
			gameFinish = 1
			continue
			//break
		}

		turn = (turn + 1) % 2
		played <- true
	}
	pconn.Close()
}
