package main

import "net"
import "fmt"
import "bufio"
import "strings" // only needed below for sample processing

func exe_cmd(cmd string) (string, error) {
  fmt.Println("command is ",cmd)
  // splitting head => g++ parts => rest of the command
  parts := strings.Fields(cmd)
  head := parts[0]
  parts = parts[1:len(parts)]

  out, err := exec.Command(head,parts...).Output()
  if err != nil {
    fmt.Printf("%s", err)
  }
  fmt.Printf("%s", out)
  return out, err 
}

func send_response_to_collector(response string) {
  conn, _ := net.Dial("tcp", "127.0.0.1:8082")
  // send to socket
  fmt.Fprintf(conn, response + "\n")
  // listen for reply
  message, _ := bufio.NewReader(conn).ReadString('\n')
  fmt.Print("Message from server: "+message)
}

func main() {

  fmt.Println("Launching server...")

  // listen on all interfaces
  ln, _ := net.Listen("tcp", ":8081")

  // accept connection on port
  conn, _ := ln.Accept()

  // run loop forever (or until ctrl-c)
  for {
    // will listen for message to process ending in newline (\n)
    message, _ := bufio.NewReader(conn).ReadString('\n')
    // output message received
    fmt.Print("Message Received:", string(message))
    // sample process for string received
    newmessage := exe_cmd(message)
    // send new string back to client
    conn.Write([]byte("ack\n"))

    send_response_to_collector(newmessage)
  }
}