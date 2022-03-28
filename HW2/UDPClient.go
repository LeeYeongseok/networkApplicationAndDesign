/**
 * 20170768 이영석
 * UDPClient.go
 **/

package main

import ("bufio"; "fmt"; "net"; "os" ; "time" ; "os/signal" ; "syscall" )

// 소수점 셋째자리까지 잘라서 문자열로 반환
func cutFloat(num float64) string{
  return fmt.Sprintf("%06.3f",num)
}

// 서버에 에러가 있으면 알리고 close 및 종료
func checkError(err error, pconn net.PacketConn){
  if nil != err {
        fmt.Printf("server Disconnected\n")
        pconn.Close()
        os.Exit(0)
      }
}

func main() {
  serverName := "nsl2.cau.ac.kr"
  serverPort := "30768"
  server_addr, _ := net.ResolveUDPAddr("udp", serverName+":"+serverPort)
  pconn, _:= net.ListenPacket("udp", ":")
  localAddr := pconn.LocalAddr().(*net.UDPAddr)
  fmt.Printf("Client is running on port %d\n", localAddr.Port)

// 키보드 인터럽트시 연결 끊고 종료
  signals := make( chan os.Signal, 1) 
  signal.Notify( signals, syscall.SIGINT, syscall.SIGTERM) 
  go func(){
    <- signals
    fmt.Println("\nBye bye~")
    pconn.Close()
    os.Exit(0) 
  }()

  
  for{
    fmt.Printf("<Menu>\n1) convert text to UPPER-case\n2) get my IP address and port number\n3) get server request count\n4) get server running time\n5) exit\nInput option: ")
    num, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    
    ch1 := make(chan bool) // 타임아웃 확인 채널
    ch2 := make(chan bool) // 통신완료 확인 채널

    buffer := make([]byte, 1024)
    var elapsedTime float64 // 경과 시간
    
  	if num == "1\n" {
      fmt.Printf("Input sentence: ")
      input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
      pconn.WriteTo([]byte("1"+ input),server_addr)
      
  	}else if num == "2\n"{
      pconn.WriteTo([]byte("2"),server_addr)

  	}else if num == "3\n"{
      pconn.WriteTo([]byte("3"),server_addr)
      
  	}else if num == "4\n"{
      pconn.WriteTo([]byte("4"),server_addr)

  	}else if num == "5\n"{ // 종료
  		break
  	}else{ // 다른문자 입력시
      continue
    }

    // 타임아웃 5초 설정
    go func(timeout chan bool){
      time.Sleep(5 * time.Second)
      timeout <- true
    }(ch1)

    // 서버와 통신
    go func(done chan bool){
      startTime := time.Now()
      pconn.ReadFrom(buffer)
      elapsedTime = float64(time.Since(startTime)) / float64(time.Millisecond)
      done <- true
    }(ch2)
  
    select{
        case <- ch1: // 타임아웃 시 메시지 출력 후 연결 끊고 종료
          fmt.Println("Time out : Cannot connect with server")
          fmt.Println("Bye bye~")
          pconn.Close()
          os.Exit(0)

        case <- ch2: // 측정한 시간 및 서버에서 전달받은 내용 출력
          fmt.Printf("\nReply from server: %s", string(buffer))
          fmt.Println("RTT = ", cutFloat(elapsedTime),"ms\n")
    }
    
  }
  
  fmt.Println("Bye bye~")
  pconn.Close()
}
