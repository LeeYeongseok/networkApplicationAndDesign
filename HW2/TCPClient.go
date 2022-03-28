/**
 * 20170768 이영석
 * TCPClient.go
 **/

package main

import ("bufio"; "fmt"; "net"; "os" ; "time" ; "os/signal" ; "syscall")

// 소수점 셋째자리까지 잘라서 문자열로 반환
func cutFloat(num float64) string{
  return fmt.Sprintf("%.3f",num)
}

// 서버에 에러가 있으면 알리고 close 및 종료
func checkError(err error, conn net.Conn){
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
  conn, err:= net.Dial("tcp", serverName+":"+serverPort)
  if nil != err {
        fmt.Println("server Disconnected")
        fmt.Println("Bye bye~")
        os.Exit(0)
      }
  localAddr := conn.LocalAddr().(*net.TCPAddr)
  fmt.Printf("Client is running on port %d\n", localAddr.Port)

// 키보드 인터럽트시 연결 끊고 종료
  signals := make( chan os.Signal, 1) 
  signal.Notify( signals, syscall.SIGINT, syscall.SIGTERM) 
  go func(){
    <- signals
    fmt.Println("\nBye bye~")
    conn.Close()
    os.Exit(0) 
  }()

  
  for{
    fmt.Printf("<Menu>\n1) convert text to UPPER-case\n2) get my IP address and port number\n3) get server request count\n4) get server running time\n5) exit\nInput option: ")
    num, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    buffer := make([]byte, 1024)
    var elapsedTime float64 // 경과시간
    
  	if num == "1\n" {
      fmt.Printf("Input sentence: ")
      input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
      conn.Write([]byte("1"+ input))
      //시간 시작
      startTime := time.Now()
      _, err := conn.Read(buffer)
      //시간 종료      
      elapsedTime = float64(time.Since(startTime)) / float64(time.Millisecond)
      checkError(err, conn)
      
  	}else if num == "2\n"{
      conn.Write([]byte("2"))
      //시간 시작
      startTime := time.Now()
      _, err := conn.Read(buffer)
      //시간 종료      
      elapsedTime = float64(time.Since(startTime)) / float64(time.Millisecond)
      checkError(err, conn)
      
  	}else if num == "3\n"{
      conn.Write([]byte("3"))
      //시간 시작
      startTime := time.Now()
      _, err := conn.Read(buffer)
      //시간 종료      
      elapsedTime = float64(time.Since(startTime)) / float64(time.Millisecond)
      checkError(err, conn)
      
  	}else if num == "4\n"{
      conn.Write([]byte("4"))
      //시간 시작
      startTime := time.Now()
      _, err := conn.Read(buffer)
      //시간 종료
      elapsedTime = float64(time.Since(startTime)) / float64(time.Millisecond)
      checkError(err, conn)
      
  	}else if num == "5\n"{
  		break
  	}else{
      continue
    }

    fmt.Printf("\nReply from server: %s", string(buffer))
    fmt.Println("RTT = ", cutFloat(elapsedTime),"ms\n") 
    
  }
  fmt.Println("Bye bye~")
  conn.Close()
}
