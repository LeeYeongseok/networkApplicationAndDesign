/**
 * 20170768 이영석
 * UDPServer.go
 **/

package main

import ("bytes"; "fmt"; "net" ; "strings" ; "time" ; "os" ; "os/signal" ; "syscall")

func main() {
  startTime := time.Now()
  serverPort := "30768"
  fmt.Printf("Server is ready to receive on port %s\n", serverPort)
  buffer := make([]byte, 1024)
  countCommand := 0
  var pconn net.PacketConn
  var err error

  signals := make( chan os.Signal, 1) 
  signal.Notify( signals, syscall.SIGINT, syscall.SIGTERM) 
  go func(){
    <- signals
    fmt.Println("\nBye bye~")
    if pconn!=nil{
     pconn.Close() 
    }
    os.Exit(0) 
  }()


  pconn, _= net.ListenPacket("udp", ":"+serverPort)
  if nil != err{
    if pconn!=nil{
     pconn.Close() 
    }
    os.Exit(0) 
  }

  for{

    count, r_addr, _:= pconn.ReadFrom(buffer)
    if nil != err{
       break
    }
    fmt.Printf("Connection request from %s\n", r_addr.String())

    if string(buffer[0:1]) == "1"{
      pconn.WriteTo(bytes.ToUpper(buffer[1:count]), r_addr)
      fmt.Printf("Command 1\n")
      countCommand += 1
    } else if string(buffer[0:1]) == "2"{
      point := strings.Index(r_addr.String(), ":")
      IP := r_addr.String()[:point]
      fmt.Printf("Command 2\n")
      port := r_addr.String()[point+1:]
      reply := "client IP = "+ IP + ", port = "+ port + "\n"
      pconn.WriteTo([]byte(reply),r_addr)
      countCommand += 1
    } else if string(buffer[0:1]) == "3"{
      pconn.WriteTo([]byte("requests served = "+fmt.Sprintf("%02d", countCommand)+"\n"),r_addr)
      fmt.Printf("Command 3\n")
      countCommand += 1
    } else if string(buffer[0:1]) == "4"{
      elapsedTime := int(float64(time.Since(startTime)) / float64(time.Second))
      min := elapsedTime / 60
      hour := min / 60
      sec := elapsedTime % 60
      min = min % 60
      nowTime := fmt.Sprintf("%02d:%02d:%02d",hour,min,sec)
      pconn.WriteTo([]byte("run time = " + nowTime +"\n"),r_addr)
      fmt.Printf("Command 4\n")
      countCommand += 1
    }
  }
  pconn.Close()
  pconn = nil
}