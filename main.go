package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	const HOST string = "localhost"
	const PORT int = 8080
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", HOST, PORT))
	if err != nil {
		fmt.Printf("Unable to connect to host %v and port %v\n", HOST, PORT)
		return
	}
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go receiveData(conn, &waitGroup)
	go sendData(conn)
	waitGroup.Wait()
}

func receiveData(conn net.Conn, waitGroup *sync.WaitGroup) {
	buffer := make([]byte, 1028)
	for {
		dataLength, err := conn.Read(buffer)
		if err != nil {
			log.Println("Disconnected from server")
			conn.Close()
			waitGroup.Done()
			break
		}
		dataBytes := buffer[:dataLength]
		data := string(dataBytes)
		log.Print(data)
	}
}

func sendData(conn net.Conn) {
	fmt.Print("Select your username: ")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		input := scanner.Text()
		dataBytes := []byte(input)
		data := string(dataBytes)
		if data == "END" {
			conn.Close()
			log.Println("Disconnected")
			break
		}
		_, err := conn.Write(dataBytes)
		if err != nil {
			log.Println(err)
			conn.Close()
			break
		}
	}
}
