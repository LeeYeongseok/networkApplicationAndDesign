/**
 * 20170768 이영석
 * TCPServer.go
 **/

package main

import ("bytes"; "fmt"; "net" ; "strings" ; "time" ; "os" ; "os/signal" ; "syscall")

func main() {
  startTime := time.Now()
  serverPort := "30768"
  listener, _:= net.Listen("tcp", ":" + serverPort)
  fmt.Printf("Server is ready to receive on port %s\n", serverPort)
  buffer := make([]byte, 1024)
  countCommand := 0
  var conn net.Conn
  var err error

  //키보드 인터럽트 시 메시지 출력 후 close 및 종료
  signals := make( chan os.Signal, 1) 
  signal.Notify( signals, syscall.SIGINT, syscall.SIGTERM) 
  go func(){
    <- signals
    fmt.Println("\nBye bye~")
    if conn!=nil{
     conn.Close() 
    }
    os.Exit(0) 
  }()

  for {
    conn, err = listener.Accept()
    if nil != err{
      break
    }
    fmt.Printf("Connection request from %s\n", conn.RemoteAddr().String())
    for{
      count, err := conn.Read(buffer)
      if nil != err{
         break
      }
      
    	if string(buffer[0:1]) == "1"{
    		conn.Write(bytes.ToUpper(buffer[1:count]))
		    fmt.Printf("Command 1\n")
        countCommand += 1
        
    	} else if string(buffer[0:1]) == "2"{
    		point := strings.Index(conn.RemoteAddr().String(), ":")
    		IP := conn.RemoteAddr().String()[:point]
		    fmt.Printf("Command 2\n")
        port := conn.RemoteAddr().String()[point+1:]
    		reply := "client IP = "+ IP + ", port = "+ port + "\n"
    		conn.Write([]byte(reply))
        countCommand += 1
        
    	} else if string(buffer[0:1]) == "3"{
        conn.Write([]byte("requests served = "+fmt.Sprintf("%02d",countCommand)+"\n"))
		    fmt.Printf("Command 3\n")
        countCommand += 1
        
      } else if string(buffer[0:1]) == "4"{
        elapsedTime := int(float64(time.Since(startTime)) / float64(time.Second))
        min := elapsedTime / 60
        hour := min / 60
        sec := elapsedTime % 60
        min = min % 60
        nowTime := fmt.Sprintf("%02d:%02d:%02d",hour,min,sec)
        conn.Write([]byte("run time = " + nowTime +"\n"))
		    fmt.Printf("Command 4\n")
        countCommand += 1
      }
    }
    conn.Close()
  }
}